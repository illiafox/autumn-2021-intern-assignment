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

func New() App {
	return App{}
}

func (app *App) Init() DeferFunc {

	var (
		HTTP       = flag.Bool("http", false, "force HTTP mode")
		logPath    = flag.String("log", "log.txt", "log file path (default 'log.txt')")
		configPath = flag.String("config", "config.toml", "config path (default 'config.toml')")

		env = flag.Bool("env", false, "load from environment variables")
	)
	flag.Parse()

	// // //

	file, err := os.OpenFile(*logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln(fmt.Errorf("creating/opening log file (%s): %w", *logPath, err))
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatalln(fmt.Errorf("getting file stats: %w", err))
	}

	if info.Size() > 0 {
		_, err = file.Write([]byte("\n\n"))
		if err != nil {
			log.Fatalln(fmt.Errorf("write to file: %w", err))
		}
	}

	logger := loggers.NewLogger(file)

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		logger.Fatal("reading config file", zap.String("config", *configPath), zap.Error(err))
	}

	if *env {
		err = conf.LoadEnv()
		if err != nil {
			logger.Fatal("loading environments", zap.Error(err))
		}
	}

	app.logger = logger
	app.conf = conf
	app.HTTP = *HTTP

	runtime.GC()

	// sync function
	return func() {
		err := logger.Sync()
		if err != nil {
			_, ok := err.(*fs.PathError)
			if ok {
				return
			}

			log.Println(fmt.Errorf("sync logger: %w", err))
		}
	}
}
