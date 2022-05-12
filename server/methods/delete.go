package methods

import (
	"autumn-2021-intern-assignment/public"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (m Methods) Delete(w http.ResponseWriter, r *http.Request) {

	var del = struct {
		User int64 `json:"user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&del)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		EncodeError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if del.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		EncodeError(w, fmt.Errorf("wrong 'user' field value: %d", del.User))

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
		EncodeError(w, fmt.Errorf("delete user: %w", err))

		return
	}

	w.Write(success)
}
