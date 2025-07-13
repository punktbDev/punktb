package service

import (
	"gitlab.com/freelance/punkt-b/backend/internal/domain"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
)

type (
	manager struct {
		dm domain.Manager
	}
	Manager interface {
		GetUserData(login, password string) (*dto.Manager, error)
		ChangeManagerData(ms *dto.Manager) error
		GetAllManagers() ([]dto.Manager, error)
		AddManager(mn dto.Manager) error
		ChangeActive(id int) error
		ChangeFullAccess(id int) error
	}
)

func NewManager(dm domain.Manager) Manager {
	return &manager{dm: dm}
}

func (m *manager) ChangeFullAccess(id int) error {
	if err := m.dm.ChangeFullAccess(id); err != nil {
		return err
	}

	return nil
}

func (m *manager) ChangeActive(id int) error {
	if err := m.dm.ChangeActive(id); err != nil {
		return err
	}

	return nil
}

func (m *manager) AddManager(mn dto.Manager) error {
	if err := m.dm.AddManager(mn); err != nil {
		return err
	}

	return nil
}

func (m *manager) GetAllManagers() ([]dto.Manager, error) {
	ms, err := m.dm.GetAllManagers()
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (m *manager) GetUserData(login, password string) (*dto.Manager, error) {
	k, err := m.dm.GetUserData(login, password)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (m *manager) ChangeManagerData(ms *dto.Manager) error {
	if err := m.dm.ChangeManagerData(ms); err != nil {
		return err
	}

	return nil
}
