package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"autumn-2021-intern-assignment/utils/config"
	// postgre
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
	for err != nil {
		return nil, fmt.Errorf("connecting to mysql: %w", err)
	}

	err = conn.Ping()

	for i := 0; err != nil && i < 5; i++ {
		log.Println(fmt.Errorf("ping: %w -> retrying", err))
		time.Sleep(time.Second)

		err = conn.Ping()
	}

	if err != nil {
		os.Exit(1)
	}

	return &DB{conn}, nil
}

func NewFromConnect(conn *sql.DB) *DB {
	return &DB{conn}
}
