package app

import (
	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/exchange"
	"go.uber.org/zap"
)

func (app *App) Database() (DeferFunc, bool) {
	app.logger.Info("Initializing database")

	db, err := database.New(app.conf.Postgres)
	if err != nil {
		app.logger.Error("connecting to database", zap.Error(err))

		return nil, false
	}

	exchange.UpdateWithLoad(app.conf.Exchanger, app.logger)

	app.db = db

	// close func
	return func() {
		app.logger.Info("Closing database connection")
		db.Close()
	}, true
}
