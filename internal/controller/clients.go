package controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/internal/service"
	"net/http"
	"strconv"
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

	if len(cls) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		SendResponse(http.StatusOK, w, cls)
	}
}
