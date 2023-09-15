package controller

import (
	"database/sql"
	"errors"
	"gitlab.com/freelance/punkt-b/backend/internal/service"
	"net/http"
)

type (
	login struct {
		srv service.Login
		err ErrorResponse
	}
	Login interface {
		Login(w http.ResponseWriter, r *http.Request)
	}
)

func NewLogin(srv service.Login) Login {
	return &login{srv: srv, err: NewError()}
}

func (l *login) Login(w http.ResponseWriter, r *http.Request) {
	m := GetManager(r)

	m, err := l.srv.Login(m.Login, m.Password)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			l.err.BadRequest(w, err, http.StatusUnauthorized)
		} else {
			l.err.BadRequest(w, err, http.StatusInternalServerError)
		}
		return
	}

	SendResponse(http.StatusOK, w, m)
}
