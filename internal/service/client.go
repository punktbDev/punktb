package service

import (
	"gitlab.com/freelance/punkt-b/backend/internal/domain"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
)

type (
	client struct {
		dm domain.Client
	}
	Client interface {
		GetClients(id int, isAdmin bool) ([]dto.Client, error)
		SetClientChecked(id int) error
		SetClientArchive(id int) error
		AddResult(cl *dto.Client) error
		GetResultClient(id int) (*dto.Client, error)
	}
)

func NewClient(dm domain.Client) Client {
	return &client{dm: dm}
}

func (c *client) GetResultClient(id int) (*dto.Client, error) {
	cl, err := c.dm.GetResultClient(id)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *client) AddResult(cl *dto.Client) error {
	if err := c.dm.AddResult(cl); err != nil {
		return err
	}

	return nil
}

func (c *client) GetClients(id int, isAdmin bool) ([]dto.Client, error) {
	m, err := c.dm.GetClients(id, isAdmin)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (c *client) SetClientChecked(id int) error {
	if err := c.dm.SetClientChecked(id); err != nil {
		return err
	}

	return nil
}

func (c *client) SetClientArchive(id int) error {
	if err := c.dm.SetClientArchive(id); err != nil {
		return err
	}

	return nil
}
