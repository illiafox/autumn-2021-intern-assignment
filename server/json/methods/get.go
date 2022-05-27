package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"autumn-2021-intern-assignment/exchange"
	"autumn-2021-intern-assignment/public"
	"github.com/shopspring/decimal"
)

func (m Methods) Get(w http.ResponseWriter, r *http.Request) {
	var get = struct {
		User int64  `json:"user_id"`
		Base string `json:"base"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&get)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		EncodeErr(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if get.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		EncodeErr(w, fmt.Errorf("wrong 'user' field value: %d", get.User))

		return
	}

	balance, _, err := m.db.GetBalance(context.Background(), get.User)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusInternalServerError)
			EncodeErr(w, fmt.Errorf("get balance: %w", err))
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			EncodeErr(w, fmt.Errorf("balance with user id %d not found", get.User))
		}

		return
	}

	var ret = struct {
		Ok      bool    `json:"ok"`
		Base    string  `json:"base"`
		Rate    float64 `json:"rate,omitempty"`
		Balance string  `json:"balance"`
	}{Ok: true}

	if get.Base != "" {
		ex, ok := exchange.GetExchange(get.Base)
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			EncodeErr(w, fmt.Errorf("base: abbreviation '%s' is not supported", get.Base))

			return
		}
		ret.Rate = ex
		ret.Balance = decimal.NewFromFloat(float64(balance) / 100).Div(decimal.NewFromFloat(ex)).StringFixed(2)
	} else {
		get.Base = "RUB"
		ret.Balance = strconv.FormatFloat(float64(balance)/100, 'f', 2, 64)
	}

	ret.Base = get.Base

	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		EncodeErr(w, fmt.Errorf("encoding json: %w", err))

		return
	}
}
