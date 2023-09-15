package domain

import (
	dbModel "gitlab.com/freelance/punkt-b/backend/internal/database"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/pkg/database"
)

type (
	client struct {
		db database.Database
	}
	Client interface {
		GetClients(id int, isAdmin bool) ([]dto.Client, error)
		SetClientChecked(id int) error
		SetClientArchive(id int) error
		AddResult(cl *dto.Client) error
		GetResultClient(id int) (*dto.Client, error)
	}
)

func NewClient(db database.Database) Client {
	return &client{db: db}
}

func (c *client) GetResultClient(id int) (*dto.Client, error) {
	cl, err := c.db.Read(&dbModel.Client{Id: id})
	if err != nil {
		return nil, err
	}

	return cl.(*dto.Client), err
}

func (c *client) AddResult(cl *dto.Client) error {
	_, err := c.db.Create(&dbModel.ClientResult{
		ManagerId: cl.ManagerId,
		Name:      cl.Name,
		Phone:     cl.Phone,
		Email:     cl.Email,
		New:       cl.New,
		InArchive: cl.InArchive,
		Results:   cl.Result,
		Date:      cl.Date,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) GetClients(id int, isAdmin bool) ([]dto.Client, error) {
	cls, err := c.db.Read(&dbModel.GetClients{
		Id:      id,
		IsAdmin: isAdmin,
	})
	if err != nil {
		return nil, err
	}

	return cls.([]dto.Client), nil
}

func (c *client) SetClientChecked(id int) error {
	if err := c.db.Update(&dbModel.Client{Id: id}); err != nil {
		return err
	}

	return nil
}

func (c *client) SetClientArchive(id int) error {
	if err := c.db.Update(&dbModel.ClientInArchive{Id: id}); err != nil {
		return err
	}

	return nil
}
