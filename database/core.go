package database

import (
	"autumn-2021-intern-assignment/utils/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"

	// postgres
	_ "github.com/jackc/pgx/v4"
)

type DB struct {
	conn *pgxpool.Pool
}

func New(conf config.Postgres) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable",
			conf.User,
			conf.Pass,
			conf.IP,
			conf.Port,
			conf.DbName,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{pool}, nil
}
