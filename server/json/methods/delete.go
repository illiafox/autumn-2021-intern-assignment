package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/public"
)

// Delete
// @Description User ID
type Delete struct {
	User int64 `json:"user_id"`
}

// Delete godoc
// @Summary      Delete Balance
// @Description  Delete User balance
// @Accept       json
// @Produce      json
// @Param        input body 	Delete true "User ID"
// @Success      200  {boolean} true
// @Failure      400  {object}  Error
// @Failure      406  {object}  Error
// @Failure      500  {object}  Error
// @Router       /delete [delete]
func (m Methods) Delete(w http.ResponseWriter, r *http.Request) {

	var del Delete

	err := json.NewDecoder(r.Body).Decode(&del)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if del.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'user' field value: %d", del.User))

		return
	}

	ctx := context.Background()

	err = m.db.Delete(ctx, del.User)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
		WriteError(w, fmt.Errorf("delete user: %w", err))

		return
	}

	w.Write(success)
}
