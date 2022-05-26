package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"autumn-2021-intern-assignment/server"
	"go.uber.org/zap"
)

func (app *App) Start() {
	srv := server.New(app.db.Methods, app.conf.Host)

	ch := make(chan os.Signal, 1)

	var err error

	go func() {
		app.logger.Info("Server started at " + srv.Addr)

		if app.HTTP {
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS(app.conf.Host.Cert, app.conf.Host.Key)
		}

		if err != nil {
			if err != http.ErrServerClosed {
				app.logger.Error("Server", zap.Error(err))
			}
			ch <- nil
		}
	}()

	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-ch
	os.Stdout.WriteString("\n")

	// Create a deadline to wait for closing all connections
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	app.logger.Info("Shutting down server")
	err = srv.Shutdown(ctx)
	if err != nil {
		app.logger.Error("Shutting:", zap.Error(err))
	}
}
