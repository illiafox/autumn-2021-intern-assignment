package database

import (
	"context"
	"fmt"
	"time"

	"autumn-2021-intern-assignment/database/methods"
	"autumn-2021-intern-assignment/utils/config"
	// postgres
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(conf config.Postgres) (*Database, error) {
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

	return &Database{
		Methods: methods.New(pool),
		pool:    pool,
	}, nil
}

type Database struct {
	Methods methods.Methods
	pool    *pgxpool.Pool
}

func (d Database) Close() {
	d.pool.Close()
}
