package methods

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/exchange"
	"autumn-2021-intern-assignment/public"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type getJSON struct {
	User int64  `json:"user_id"`
	Base string `json:"base"`
}

type getRetJSON struct {
	Ok      bool   `json:"ok"`
	Base    string `json:"base"`
	Balance string `json:"balance"`
}

func Get(db *database.DB, ctx *fasthttp.RequestCtx) {
	var get getJSON

	err := json.Unmarshal(ctx.PostBody(), &get)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "decoding json: %w", err)

		return
	}

	if get.User <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "wrong 'user' field value: %d", get.User)

		return
	}

	balance, _, err := db.GetBalance(context.Background(), get.User)
	if err != nil {
		if errors.As(err, public.ErrInternal) {
			ctx.SetStatusCode(http.StatusInternalServerError)
		} else {
			ctx.SetStatusCode(http.StatusNotAcceptable)
		}
		jsonError(ctx, "get balance: %w", err)

		return
	}

	var ret = getRetJSON{Ok: true}

	if get.Base != "" {
		ex, ok := exchange.GetExchange(get.Base)
		if !ok {
			ctx.SetStatusCode(http.StatusNotAcceptable)
			jsonError(ctx, "base: abbreviation '%s' is not supported", get.Base)

			return
		}

		ret.Balance = decimal.NewFromFloat(float64(balance) / 100).Div(decimal.NewFromFloat(ex)).StringFixed(2)
	} else {
		get.Base = "RUB"
		ret.Balance = strconv.FormatFloat(float64(balance)/100, 'f', 2, 64)
	}

	ret.Base = get.Base

	data, err := json.Marshal(ret)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		jsonError(ctx, "encoding json: %w", err)

		return
	}

	ctx.Write(data)
}
