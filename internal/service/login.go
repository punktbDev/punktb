package service

import (
	"gitlab.com/freelance/punkt-b/backend/internal/domain"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
)

type (
	login struct {
		dm domain.Login
	}
	Login interface {
		Login(login, password string) (*dto.Manager, error)
	}
)

func NewLogin(dm domain.Login) Login {
	return &login{dm: dm}
}

func (l *login) Login(login, password string) (*dto.Manager, error) {
	m, err := l.dm.Login(login, password)
	if err != nil {
		return nil, err
	}

	return m, nil
}
