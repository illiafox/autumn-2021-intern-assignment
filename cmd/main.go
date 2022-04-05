package main

import (
	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/exchange"
	"autumn-2021-intern-assignment/server"
	"autumn-2021-intern-assignment/utils/config"
	"autumn-2021-intern-assignment/utils/multiwriter"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "config.toml", "config path (default 'config.toml')")
	logPath := flag.String("log", "log.txt", "log file path (default 'log.txt')")
	flag.Parse()

	{
		l, err := os.Create(*logPath)
		if err != nil {
			log.Fatalln(fmt.Errorf("creating log file (%s): %w", *logPath, err))
		}
		log.SetOutput(multiwriter.New(os.Stdout, l))
	}

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("reading config file (%s): %w", *configPath, err))
	}

	db, err := database.New(conf.MySQL)
	if err != nil {
		log.Fatalln(fmt.Errorf("connecting to database: %w", err))
	}

	if !conf.Exchanger.Skip {
		err = exchange.Update(conf.Exchanger)
		if err != nil {
			log.Fatalln(fmt.Errorf("updating currencies: %w", err))
		}

		time.Sleep(time.Second * time.Duration(conf.Exchanger.Every))
		go func() {
			time.Sleep(time.Second * time.Duration(conf.Exchanger.Every))
			err = exchange.Update(conf.Exchanger)
			for err != nil {
				log.Println(fmt.Errorf("updating currencies: %w", err))
				err = exchange.Update(conf.Exchanger)
			}
		}()
	}

	server.Start(db, conf.Host)
}
