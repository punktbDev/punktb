package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrUserExist = errors.New("пользователь с таким логином уже существует")
)

type (
	manager struct {
		srv service.Manager
		err errorResponse
	}
	Manager interface {
		SelfManager(w http.ResponseWriter, r *http.Request)
		GetUserData(w http.ResponseWriter, r *http.Request)
		ChangeManagerData(w http.ResponseWriter, r *http.Request)
		GetAllManagers(w http.ResponseWriter, r *http.Request)
		AddManager(w http.ResponseWriter, r *http.Request)
		ChangeActive(w http.ResponseWriter, r *http.Request)
		ChangeFullAccess(w http.ResponseWriter, r *http.Request)
	}
)

func NewManager(srv service.Manager) Manager {
	return &manager{srv: srv}
}

func (m *manager) ChangeFullAccess(w http.ResponseWriter, r *http.Request) {
	mn := GetManager(r)
	if !mn.IsAdmin {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	tId, ok := vars["id"]
	if !ok {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(tId)
	if err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err = m.srv.ChangeFullAccess(id); err != nil {
		if errors.Is(err, database.ErrManagerNotFound) {
			m.err.BadRequest(w, err, http.StatusBadRequest)
		} else {
			m.err.BadRequest(w, err, http.StatusInternalServerError)
		}

		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (m *manager) ChangeActive(w http.ResponseWriter, r *http.Request) {
	mn := GetManager(r)
	if !mn.IsAdmin {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	tId, ok := vars["id"]
	if !ok {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(tId)
	if err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err = m.srv.ChangeActive(id); err != nil {
		if errors.Is(err, database.ErrDeleteAdmin) {
			m.err.BadRequest(w, err, http.StatusForbidden)
		} else {
			m.err.BadRequest(w, err, http.StatusInternalServerError)
		}
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (m *manager) AddManager(w http.ResponseWriter, r *http.Request) {
	ms := GetManager(r)
	if !ms.IsAdmin {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var mn dto.Manager
	if err := decoder.Decode(&mn); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := mn.ValidateCreate(); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := m.srv.AddManager(mn); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			m.err.BadRequest(w, ErrUserExist, http.StatusBadRequest)
		} else {
			m.err.BadRequest(w, err, http.StatusInternalServerError)
		}
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (m *manager) GetAllManagers(w http.ResponseWriter, r *http.Request) {
	mn := GetManager(r)
	if !mn.IsAdmin {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}

	ms, err := m.srv.GetAllManagers()
	if err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	SendResponse(http.StatusOK, w, &dto.ResponseData{Data: ms})
}

func (m *manager) GetUserData(w http.ResponseWriter, r *http.Request) {
	mn := GetManager(r)
	ms, err := m.srv.GetUserData(mn.Login, mn.Password)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			m.err.BadRequest(w, err, http.StatusUnauthorized)
		} else {
			m.err.BadRequest(w, err, http.StatusInternalServerError)
		}

		return
	}

	SendResponse(http.StatusOK, w, ms)
}

func (m *manager) ChangeManagerData(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	vars := mux.Vars(r)

	i := vars["id"]
	var mg dto.Manager
	if i == "" {
		m.err.BadRequest(w, errors.New("manager id is missing"), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return

	}

	if err = decoder.Decode(&mg); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if !GetManager(r).IsAdmin {
		m.err.BadRequest(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}

	mg.Id = id

	if err = mg.ValidateUpdate(); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err = m.srv.ChangeManagerData(&mg); err != nil {
		m.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}

func (m *manager) SelfManager(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	mg := GetManager(r)
	zap.L().Info("change data for manager",
		zap.String("name", mg.Name),
		zap.String("surname", mg.Surname),
		zap.String("login", mg.Login))

	var income dto.Manager
	if err := decoder.Decode(&income); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	income.Id = mg.Id

	if err := income.ValidateForSelfUpdate(); err != nil {
		m.err.BadRequest(w, err, http.StatusBadRequest)
		return
	}

	if err := m.srv.SelfManager(&income); err != nil {
		m.err.BadRequest(w, err, http.StatusInternalServerError)
		return
	}

	SendResponse(http.StatusOK, w, &dto.Response{Success: true})
}
