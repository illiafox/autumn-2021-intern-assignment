package database

import (
	"autumn-2021-intern-assignment/utils/config"
	"database/sql"
	"fmt"
	// postgres
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func New(conf config.Postgres) (*DB, error) {
	conn, err := sql.Open(
		"postgres",
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
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{conn}, nil
}
