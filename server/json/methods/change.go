package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/public"
)

func (m Methods) Change(w http.ResponseWriter, r *http.Request) {

	var change = struct {
		User        int64  `json:"user_id"`
		Change      int64  `json:"change"`
		Description string `json:"description"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&change)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		EncodeErr(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if change.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		EncodeErr(w, fmt.Errorf("wrong 'user' field value: %d", change.User))

		return
	}

	if change.Change == 0 {
		w.WriteHeader(http.StatusBadRequest)
		EncodeString(w, "wrong 'change' field value: can't be zero")

		return
	}

	if change.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		EncodeString(w, "'description' field value can't be empty")

		return
	}

	ctx := context.Background()

	err = m.db.ChangeBalance(ctx, change.User, change.Change, change.Description)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}

		EncodeErr(w, fmt.Errorf("change balance: %w", err))

		return
	}

	w.Write(success)
}
