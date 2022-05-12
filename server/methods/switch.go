package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/public"
)

func (m Methods) Switch(w http.ResponseWriter, r *http.Request) {
	var sw = struct {
		OldUserID int64 `json:"old_user_id"`
		NewUserID int64 `json:"new_user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&sw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		EncodeError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if sw.OldUserID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		EncodeError(w, fmt.Errorf("'old_user_id' can't be lower or equal zero, got %d", sw.OldUserID))

		return
	}

	if sw.NewUserID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		EncodeError(w, fmt.Errorf("'new_user_id' can't be lower or equal zero, got %d", sw.NewUserID))

		return
	}

	ctx := context.Background()

	err = m.db.Switch(ctx, sw.OldUserID, sw.NewUserID)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
		EncodeError(w, fmt.Errorf("switch: %w", err))

		return
	}

	w.Write(success)
}
