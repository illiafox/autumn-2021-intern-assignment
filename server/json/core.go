package json

import (
	"net/http"

	"autumn-2021-intern-assignment/database/model"
	"autumn-2021-intern-assignment/server/json/methods"
	"autumn-2021-intern-assignment/server/json/middleware"
	"github.com/gorilla/mux"
)

func New(db model.Repository) http.Handler {

	m := methods.New(db)

	router := mux.NewRouter()
	router.Use(middleware.JSON)

	router.HandleFunc("/get", m.Get).Methods(http.MethodGet)
	router.HandleFunc("/change", m.Change).Methods(http.MethodPost)
	router.HandleFunc("/transfer", m.Transfer).Methods(http.MethodPost)
	router.HandleFunc("/view", m.View).Methods(http.MethodGet)
	router.HandleFunc("/switch", m.Switch).Methods(http.MethodPut)
	router.HandleFunc("/delete", m.Delete).Methods(http.MethodDelete)

	return router
}
