package database

import (
	"autumn-2021-intern-assignment/utils/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	//nolint:revive
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func New(conf config.MySQL) (*DB, error) {
	conn, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Login,
			conf.Pass,
			conf.Protocol,
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
