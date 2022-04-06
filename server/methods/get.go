package methods

import (
	"encoding/json"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/exchange"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"strconv"
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
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	balance, _, err := db.GetBalance(get.User)
	if err != nil {
		ctx.Write(jsonError("get balance: %w", err))
		return
	}

	if get.User <= 0 {
		ctx.Write(jsonError("wrong 'user' field value: %d", get.User))
		return
	}

	var ret = getRetJSON{Ok: true}

	if get.Base != "" {
		ex, ok := exchange.GetExchange(get.Base)
		if !ok {
			ctx.Write(jsonError("base: abbreviation '%s' is not supported", get.Base))
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
		ctx.Write(jsonError("encoding json: %w", err))
		return
	}
	ctx.Write(data)
}
