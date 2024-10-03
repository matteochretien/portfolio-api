package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type Client struct {
	Client *pgxpool.Pool
}

func NewClient() (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, err
	}

	config.MinConns = 5
	config.MaxConns = 10

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	conn, err := connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{Client: connPool}, nil
}

func (r *Client) Close() error {
	r.Client.Close()
	return nil
}
