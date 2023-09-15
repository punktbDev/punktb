package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	timeout = 3 * time.Second
)

type (
	Database interface {
		Create(i Creator) (int, error)
		Read(i Reader) (interface{}, error)
		Update(i Updater) error
		Delete(i Deleter) error
	}
	Creator interface {
		Create(ctx context.Context, conn *pgxpool.Conn) (int, error)
	}
	Reader interface {
		Read(ctx context.Context, conn *pgxpool.Conn) (interface{}, error)
	}
	Updater interface {
		Update(ctx context.Context, conn *pgxpool.Conn) error
	}
	Deleter interface {
		Delete(ctx context.Context, conn *pgxpool.Conn) error
	}
)

type postgres struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func NewDatabase(dbURL string) (Database, error) {
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return &postgres{
		pool:    pool,
		timeout: timeout,
	}, nil
}

func (p *postgres) Read(i Reader) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)

	defer cancel()

	conn, err := p.pool.Acquire(ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return i.Read(ctx, conn)
}

func (p *postgres) Create(i Creator) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	conn, err := p.pool.Acquire(ctx)

	if err != nil {
		return 0, err
	}

	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return 0, err
	}

	return i.Create(ctx, conn)
}

func (p *postgres) Update(i Updater) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	conn, err := p.pool.Acquire(ctx)

	if err != nil {
		return err
	}

	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return err
	}

	return i.Update(ctx, conn)
}

func (p *postgres) Delete(i Deleter) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	conn, err := p.pool.Acquire(ctx)

	if err != nil {
		return err
	}

	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return err
	}

	return i.Delete(ctx, conn)
}
