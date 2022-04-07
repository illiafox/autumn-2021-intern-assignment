package server

import (
	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/server/methods"
	routing "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func Handler(db *database.DB) fasthttp.RequestHandler {
	router := routing.New()
	w := Wrapper{db}

	router.Handle("GET", "/get", w.Wrap(methods.Get))
	router.Handle("POST", "/change", w.Wrap(methods.Change))
	router.Handle("POST", "/transfer", w.Wrap(methods.Transfer))
	router.Handle("GET", "/view", w.Wrap(methods.View))
	router.Handle("POST", "/delete", w.Wrap(methods.Delete))
	router.Handle("POST", "/switch", w.Wrap(methods.Switch))

	return router.Handler
}

type Wrapper struct {
	*database.DB
}

func (w Wrapper) Wrap(f func(*database.DB, *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		f(w.DB, ctx)
	}
}
