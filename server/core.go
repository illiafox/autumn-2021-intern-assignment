package server

import (
	"net/http"
	"time"

	"autumn-2021-intern-assignment/database/model"
	"autumn-2021-intern-assignment/server/methods"
	"autumn-2021-intern-assignment/utils/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(db model.Repository, conf config.Host) *http.Server {
	router := http.NewServeMux()

	m := methods.New(db)

	router.HandleFunc("/get", m.Get)
	router.HandleFunc("/change", m.Change)
	router.HandleFunc("/transfer", m.Transfer)
	router.HandleFunc("/view", m.View)
	router.HandleFunc("/switch", m.Switch)

	router.Handle("/metrics", promhttp.Handler())
	//

	return &http.Server{
		Addr: "0.0.0.0:" + conf.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
}
