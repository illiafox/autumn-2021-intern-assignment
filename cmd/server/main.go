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
)

func main() {
	configPath := flag.String("config", "config.toml", "config path (default 'config.toml')")
	logPath := flag.String("log", "log.txt", "log file path (default 'log.txt')")
	currPath := flag.String("curr", "currencies.json", "currencies file path (default 'currencies.json')")
	load := flag.Bool("load", false, "skip api loading, read currencies file")
	skip := flag.Bool("skip", false, "disable updating cycle")

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

	db, err := database.New(conf.Postgres)
	if err != nil {
		log.Fatalln(fmt.Errorf("connecting to database: %w", err))
	}

	if !conf.Exchanger.Skip {
		exchange.UpdateWithLoad(conf.Exchanger, *currPath, *skip, *load)
	}

	server.Start(db, conf.Host)
}
