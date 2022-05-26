package app

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/utils/config"
	loggers "autumn-2021-intern-assignment/utils/zap"
	"go.uber.org/zap"
)

type DeferFunc func()

type App struct {
	logger *zap.Logger
	conf   *config.Config
	//
	HTTP bool
	//
	db *database.Database
}

func Init() (App, DeferFunc) {

	var (
		HTTP    = flag.Bool("http", false, "force HTTP mode")
		logs    = flag.String("log", "log.txt", "log file path (default 'log.txt')")
		configs = flag.String("config", "config.toml", "config path (default 'config.toml')")

		env = flag.Bool("env", false, "load from environment variables")
	)
	flag.Parse()

	// // //

	file, err := os.OpenFile(*logs, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln(fmt.Errorf("create/open log file (%s): %w", *logs, err))
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatalln(fmt.Errorf("get file stats: %w", err))
	}

	if info.Size() > 0 {
		_, err = file.Write([]byte("\n\n"))
		if err != nil {
			log.Fatalln(fmt.Errorf("write to file: %w", err))
		}
	}

	logger := loggers.NewLogger(file)

	conf, err := config.ReadConfig(*configs)
	if err != nil {
		logger.Fatal("read config file", zap.String("config", *configs), zap.Error(err))
	}

	if *env {
		err = conf.LoadEnv()
		if err != nil {
			logger.Fatal("load environments", zap.Error(err))
		}
	}

	runtime.GC()

	// sync function
	return App{
		logger: logger,
		conf:   conf,
		HTTP:   *HTTP,
	}, sync(logger)
}

func sync(logger *zap.Logger) DeferFunc {
	return func() {
		err := logger.Sync()
		if err != nil {
			if _, ok := err.(*fs.PathError); ok {
				return
			}

			log.Println(fmt.Errorf("sync logger: %w", err))
		}
	}
}
