package controller

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/internal/service"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type (
	client struct {
		srv service.Client
		err ErrorResponse
	}
	Client interface {
		GetClients(w http.ResponseWriter, r *http.Request)
		SetClientChecked(w http.ResponseWriter, r *http.Request)
		SetClientArchive(w http.ResponseWriter, r *http.Request)
		AddResult(w http.ResponseWriter, r *http.Request)
		GetResultClient(w http.ResponseWriter, r *http.Request)
	}
)

func NewClient(srv service.Client) Client {
	return &client{srv: srv, err: NewError()}
}

func (c *client) GetResultClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tId, ok := vars["id"]
	if !ok {
		c.err.BadRequest(w, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(tId)
	if err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	cl, err := c.srv.GetResultClient(id)
	if err != nil {
		if errors.Is(err, database.ErrClientNotFound) {
			c.err.BadRequest(w, err, http.StatusBadRequest)
		} else {
			c.err.BadRequest(w, err, http.StatusInternalServerError)
		}
		return
	}

	SendResponse(http.StatusOK, w, cl)
}

func (c *client) AddResult(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var cl dto.Client
	if err := decoder.Decode(&cl); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := cl.ValidateAddResult(); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := c.srv.AddResult(&cl); err != nil {
		c.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (c *client) SetClientArchive(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var cl dto.Client
	if err := decoder.Decode(&cl); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := cl.Validate(); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := c.srv.SetClientArchive(cl.Id); err != nil {
		c.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (c *client) SetClientChecked(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var cl dto.Client
	if err := decoder.Decode(&cl); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := cl.Validate(); err != nil {
		c.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := c.srv.SetClientChecked(cl.Id); err != nil {
		c.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (c *client) GetClients(w http.ResponseWriter, r *http.Request) {
	m := GetManager(r)
	cls, err := c.srv.GetClients(m.Id, m.IsAdmin)
	if err != nil {
		c.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "application/json")

	if len(cls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		sendChunked(w, cls)
		//sendCompressed(w, cls)
	}

}

func sendCompressed(w http.ResponseWriter, cls []dto.Client) {
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close() // Закрытие потока gzip в конце

	w.WriteHeader(http.StatusOK)

	b, err := json.Marshal(cls)
	if err != nil {
		zap.L().Error("marshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = io.WriteString(gz, string(b)); err != nil {
		zap.L().Error("write", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(gz, "Sending large data more efficiently!\n")

	if err = json.NewEncoder(w).Encode(b); err != nil {
		zap.L().Error("Encode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SendResponse(http.StatusOK, w, cls)
}

func sendChunked(w http.ResponseWriter, cls []dto.Client) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintf(w, "["); err != nil {
		zap.L().Error("write start", zap.Error(err))
		return
	}
	flusher.Flush()

	for i := 0; i < len(cls); i += 100 {
		end := i + 100
		if end > len(cls) {
			end = len(cls)
		}

		b, err := json.Marshal(cls[i:end])
		if err != nil {
			zap.L().Error("marshal", zap.Error(err))
			continue
		}

		if i > 0 {
			if _, err = fmt.Fprintf(w, ","); err != nil {
				zap.L().Error("write comma", zap.Error(err))
				continue
			}
			flusher.Flush()
		}

		if _, err = fmt.Fprintf(w, "%s", b); err != nil {
			zap.L().Error("write chunk", zap.Error(err))
			continue
		}

		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}

	if _, err := fmt.Fprintf(w, "]"); err != nil {
		zap.L().Error("write end", zap.Error(err))
		return
	}
	flusher.Flush()
}
