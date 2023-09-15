package domain

import (
	dbModel "gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/pkg/database"
)

type (
	login struct {
		db database.Database
	}
	Login interface {
		Login(login, password string) (*dto.Manager, error)
	}
)

func NewLogin(db database.Database) Login {
	return &login{db: db}
}

func (l *login) Login(login, password string) (*dto.Manager, error) {
	user, err := l.db.Read(&dbModel.Login{Login: login, Password: password})
	if err != nil {
		return nil, err
	}

	return user.(*dto.Manager), nil
}
