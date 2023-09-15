package database

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
)

type (
	Login struct {
		Login    string
		Password string
	}
)

func (l *Login) Read(ctx context.Context, conn *pgxpool.Conn) (interface{}, error) {
	var m dto.Manager
	var p sql.NullString
	if err := conn.QueryRow(ctx, "SELECT id, name, surname, phone, is_admin, full_access, active FROM managers WHERE login = $1 AND password = $2",
		l.Login, l.Password).Scan(&m.Id, &m.Name, &m.Surname, &p, &m.IsAdmin, &m.IsFullAccess, &m.IsActive); err != nil {
		return nil, err
	}

	if p.Valid {
		m.Phone = p.String
	}

	return &m, nil
}
