package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
)

var (
	ErrManagerNotFound = errors.New("db: manager not found")
	ErrDeleteAdmin     = errors.New("db: is not allowed deleting admin")
)

type (
	Manager struct {
		Id       int    `json:"id"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Phone    string `json:"phone"`
		IsAdmin  bool   `json:"is_admin"`
		Secret   string `json:"-"`
	}
	ActiveManager struct {
		Id int
	}
	FullAccessManager struct {
		Id int
	}
)

func (f *FullAccessManager) Update(ctx context.Context, conn *pgxpool.Conn) error {
	tag, err := conn.Exec(ctx, `UPDATE managers SET full_access = NOT (SELECT full_access FROM managers WHERE id = $1) WHERE id = $2`, f.Id, f.Id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrManagerNotFound
	}

	return nil
}

func (m *ActiveManager) Update(ctx context.Context, conn *pgxpool.Conn) error {
	var isAdmin bool
	if err := conn.QueryRow(ctx, `SELECT is_admin FROM managers WHERE id = $1`, m.Id).Scan(&isAdmin); err != nil {
		return err
	}

	if isAdmin {
		return ErrDeleteAdmin
	}

	tag, err := conn.Exec(ctx, `UPDATE managers SET active = NOT (SELECT active FROM managers WHERE id = $1) WHERE id = $2`, m.Id, m.Id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrManagerNotFound
	}

	return nil
}

func (m *Manager) Create(ctx context.Context, conn *pgxpool.Conn) (int, error) {
	_, err := conn.Exec(ctx, `INSERT INTO managers (login, password, name, surname, phone, active, is_admin) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		m.Login, m.Password, m.Name, m.Surname, m.Phone, true, false)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (m *Manager) Update(ctx context.Context, conn *pgxpool.Conn) error {
	tag, err := conn.Exec(ctx, `UPDATE managers SET login=$1, password=$2, name=$3, surname=$4, phone=$5 WHERE id = $6`,
		m.Login, m.Password, m.Name, m.Surname, m.Phone, m.Id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrManagerNotFound
	}
	return nil
}

func (m *Manager) Read(ctx context.Context, conn *pgxpool.Conn) (interface{}, error) {
	var ms []dto.Manager
	rows, err := conn.Query(ctx, `SELECT id, login, password, name, surname, phone, active, is_admin, full_access  FROM managers`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var mn dto.Manager
		var p sql.NullString
		var a sql.NullBool

		if err = rows.Scan(&mn.Id, &mn.Login, &mn.Password, &mn.Name, &mn.Surname, &p, &a, &mn.IsAdmin, &mn.IsFullAccess); err != nil {
			return nil, err
		}

		if p.Valid {
			mn.Phone = p.String
		}

		if a.Valid {
			mn.IsActive = a.Bool
		}
		ms = append(ms, mn)
	}

	return ms, nil
}
