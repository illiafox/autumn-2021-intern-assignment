package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/public"
)

// Transfer
// @Description To/From User ID,amount and description
type Transfer struct {
	ToID        int64  `json:"to_id"`
	FromID      int64  `json:"from_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// Transfer godoc
// @Summary      Transfer money
// @Description  Transfer money between balances
// @Accept       json
// @Produce      json
// @Param        input body 	Transfer true "Transfer"
// @Success      200  {boolean} true
// @Failure      400  {object}  Error
// @Failure      406  {object}  Error
// @Failure      422  {object}  Error
// @Router       /transfer [post]
func (m Methods) Transfer(w http.ResponseWriter, r *http.Request) {
	var trans Transfer

	err := json.NewDecoder(r.Body).Decode(&trans)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if trans.ToID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'to_id' field value: %d", trans.ToID))

		return
	}
	if trans.FromID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'from_id' field value: %d", trans.FromID))

		return
	}
	if trans.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'amount' field value: can't be lower or equal zero, got %d", trans.Amount))

		return
	}

	ctx := context.Background()

	err = m.db.Transfer(ctx, trans.FromID, trans.ToID, trans.Amount, trans.Description)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
		WriteError(w, fmt.Errorf("transfer: %w", err))

		return
	}

	w.Write(success)
}
