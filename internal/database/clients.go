package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"strings"
)

var (
	ErrClientNotFound = errors.New("db: client not found")
)

type (
	GetClients struct {
		Id      int
		IsAdmin bool
	}
	Client struct {
		Id int
	}
	ClientInArchive struct {
		Id int
	}
	ClientResult struct {
		ManagerId         int
		Name              string
		Phone             string
		Email             string
		New               bool
		InArchive         bool
		Results           *dto.Result
		Date              int
		IsPhoneAdult      *bool
		ContactPermission bool
	}
)

func (c *ClientResult) Create(ctx context.Context, conn *pgxpool.Conn) (int, error) {
	var b bool
	err := conn.QueryRow(ctx, `SELECT EXISTS (SELECT email FROM clients WHERE email=$1)`, c.Email).Scan(&b)
	if err != nil {
		return 0, err
	}

	if b {
		var results []*dto.Result
		err = conn.QueryRow(ctx, `SELECT results FROM clients WHERE email=$1`, c.Email).Scan(&results)
		if err != nil {
			return 0, err
		}

		results = append(results, c.Results)
		_, err = conn.Exec(ctx, `UPDATE clients SET name=$1, phone=$2, new=$3, in_archive=$4, date=$5, results=$6, is_phone_adult = COALESCE($7, is_phone_adult), contact_permission=$8 WHERE email=$9`,
			c.Name, c.Phone, c.New, c.InArchive, c.Date, results, c.IsPhoneAdult, c.ContactPermission, c.Email)
		if err != nil {
			return 0, err
		}
	} else {
		var results []*dto.Result
		results = append(results, c.Results)
		_, err = conn.Exec(ctx, `INSERT INTO clients (manager_id, name, phone, email, new, in_archive, date, results, is_phone_adult, contact_permission) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			c.ManagerId, c.Name, c.Phone, c.Email, c.New, c.InArchive, c.Date, results, false, true)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func (c *ClientInArchive) Update(ctx context.Context, conn *pgxpool.Conn) error {
	var inArchive bool
	if err := conn.QueryRow(ctx, `SELECT in_archive FROM clients WHERE id = $1`, c.Id).Scan(&inArchive); err != nil {
		return err
	}

	_, err := conn.Exec(ctx, `UPDATE clients SET in_archive = $1 WHERE id = $2`, !inArchive, c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Update(ctx context.Context, conn *pgxpool.Conn) error {
	tag, err := conn.Exec(ctx, `UPDATE clients SET new = false WHERE id = $1`, c.Id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrClientNotFound
	}

	return nil
}

func (c *Client) Read(ctx context.Context, conn *pgxpool.Conn) (interface{}, error) {
	var cl dto.Client
	if err := conn.QueryRow(ctx, `SELECT id, manager_id, name, email, phone, new, in_archive, results, date, is_phone_adult, contact_permission FROM clients WHERE id=$1`,
		c.Id).Scan(&cl.Id, &cl.ManagerId, &cl.Name, &cl.Email, &cl.Phone, &cl.New, &cl.InArchive, &cl.Results, &cl.Date, &cl.IsPhoneAdult, &cl.ContactPermission); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return &cl, nil
}

func (g *GetClients) Read(ctx context.Context, conn *pgxpool.Conn) (interface{}, error) {
	var (
		cls  []dto.Client
		rows pgx.Rows
		err  error
	)
	if g.IsAdmin {
		rows, err = conn.Query(ctx, `SELECT * FROM clients`)
	} else {
		rows, err = conn.Query(ctx, `SELECT * FROM clients WHERE manager_id=$1`, g.Id)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c dto.Client
		if err = rows.Scan(&c.Id, &c.ManagerId, &c.Name, &c.Phone, &c.Email, &c.New, &c.InArchive, &c.Date, &c.Results, &c.IsPhoneAdult, &c.ContactPermission); err != nil {
			return nil, err
		}

		cls = append(cls, c)
	}

	return cls, nil
}
