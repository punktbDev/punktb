package domain

import (
	dbModel "gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/pkg/database"
)

type (
	manager struct {
		db database.Database
	}
	Manager interface {
		GetUserData(login, password string) (*dto.Manager, error)
		ChangeManagerData(m *dto.Manager) error
		GetAllManagers() ([]dto.Manager, error)
		AddManager(mn dto.Manager) error
		ChangeActive(id int) error
		ChangeFullAccess(id int) error
	}
)

func NewManager(db database.Database) Manager {
	return &manager{db: db}
}

func (m *manager) AddManager(mn dto.Manager) error {
	_, err := m.db.Create(&dbModel.Manager{
		Login:    mn.Login,
		Password: mn.Password,
		Name:     mn.Name,
		Surname:  mn.Surname,
		Phone:    mn.Phone,
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *manager) ChangeFullAccess(id int) error {
	if err := m.db.Update(&dbModel.FullAccessManager{Id: id}); err != nil {
		return err
	}

	return nil
}

func (m *manager) ChangeActive(id int) error {
	if err := m.db.Update(&dbModel.ActiveManager{Id: id}); err != nil {
		return err
	}

	return nil
}

func (m *manager) GetAllManagers() ([]dto.Manager, error) {
	ms, err := m.db.Read(&dbModel.Manager{})
	if err != nil {
		return nil, err
	}

	return ms.([]dto.Manager), nil
}

func (m *manager) GetUserData(login, password string) (*dto.Manager, error) {
	user, err := m.db.Read(&dbModel.Login{Login: login, Password: password})
	if err != nil {
		return nil, err
	}

	return user.(*dto.Manager), nil
}

func (m *manager) ChangeManagerData(ms *dto.Manager) error {
	if err := m.db.Update(&dbModel.Manager{
		Id:       ms.Id,
		Login:    ms.Login,
		Password: ms.Password,
		Name:     ms.Name,
		Surname:  ms.Surname,
		Phone:    ms.Phone,
	}); err != nil {
		return err
	}

	return nil
}
