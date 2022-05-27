package middleware

import "net/http"

type JSONHandler struct {
	handler http.Handler
}

func (j JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j.handler.ServeHTTP(w, r)
}

func JSON(handler http.Handler) http.Handler {
	return JSONHandler{handler: handler}
}
