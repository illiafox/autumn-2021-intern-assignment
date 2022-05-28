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

// Get
// @Description User ID and currency Base
type Get struct {
	User int64  `json:"user_id"`
	Base string `json:"base"`
}

// Balance
// @Description User Balance, Base and its Rate
type Balance struct {
	Ok bool `json:"ok"`

	Base    string `json:"base"`
	Rate    string `json:"rate,omitempty"`
	Balance string `json:"balance"`
}

// Get godoc
// @Summary      Get Balance
// @Description  Get user balance
// @Accept       json
// @Produce      json
// @Param        input body 	Get true "User id and Currency"
// @Success      200  {object}  Balance
// @Failure      400  {object}  Error
// @Failure      406  {object}  Error
// @Failure      500  {object}  Error
// @Router       /get [get]
func (m Methods) Get(w http.ResponseWriter, r *http.Request) {
	var get Get

	err := json.NewDecoder(r.Body).Decode(&get)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if get.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'user' field value: %d", get.User))

		return
	}

	balance, _, err := m.db.GetBalance(context.Background(), get.User)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusInternalServerError)
			WriteError(w, fmt.Errorf("get balance: %w", err))
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			WriteError(w, fmt.Errorf("balance with user id %d not found", get.User))
		}

		return
	}

	var ret Balance

	if get.Base != "" {
		ex, ok := exchange.GetExchange(get.Base)
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			WriteError(w, fmt.Errorf("base: abbreviation '%s' is not supported", get.Base))

			return
		}
		ret.Rate = strconv.FormatFloat(ex, 'f', 2, 64)
		ret.Balance = decimal.NewFromFloat(float64(balance) / 100).Div(decimal.NewFromFloat(ex)).StringFixed(2)
	} else {
		get.Base = "RUB"
		ret.Balance = strconv.FormatFloat(float64(balance)/100, 'f', 2, 64)
	}

	ret.Base = get.Base
	ret.Ok = true
	json.NewEncoder(w).Encode(ret)
}
