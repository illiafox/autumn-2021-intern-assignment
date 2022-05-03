package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/exchange"
	"autumn-2021-intern-assignment/server"
	"autumn-2021-intern-assignment/utils/config"
	zaplog "autumn-2021-intern-assignment/utils/zap"
	"go.uber.org/zap"
)

func main() {

	HTTP := flag.Bool("http", false, "force forceHTTP mode")

	logger, conf := Parse()
	defer logger.Sync()

	// // //

	logger.Info("Initializing database")

	db, err := database.New(conf.Postgres)
	if err != nil {
		logger.Error("connecting to database", zap.Error(err))

		return
	}
	defer func() {
		logger.Info("Closing database connection")
		db.Close()
	}()

	exchange.UpdateWithLoad(conf.Exchanger, logger)

	// // //

	srv := server.New(db.Methods, conf.Host)

	ch := make(chan os.Signal, 1)

	go func() {
		logger.Info("Server started at " + srv.Addr)

		if *HTTP {
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS(conf.Host.Cert, conf.Host.Key)
		}

		if err != nil {
			if err != http.ErrServerClosed {
				logger.Error("Server", zap.Error(err))
			}
			ch <- nil
		}
	}()

	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)

	<-ch

	// Create a deadline to wait for closing all connections
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	logger.Info("Shutting down server")
	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Shutting:", zap.Error(err))
	}

}

func Parse() (*zap.Logger, *config.Config) {
	var (
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
		file.Write([]byte("\n\n"))
	}

	logger := zaplog.NewLogger(file)

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

	return logger, conf
}
