package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/public"
)

// Switch
// @Description Old and New User ID
type Switch struct {
	OldUserID int64 `json:"old_user_id"`
	NewUserID int64 `json:"new_user_id"`
}

// Switch godoc
// @Summary      Switch owner
// @Description  Change balance owner
// @Accept       json
// @Produce      json
// @Param        input body 	Switch true "Old and New User ID"
// @Success      200  {boolean} true
// @Failure      400  {object}  Error
// @Failure      406  {object}  Error
// @Failure      422  {object}  Error
// @Router       /switch [put]
func (m Methods) Switch(w http.ResponseWriter, r *http.Request) {
	var sw Switch

	err := json.NewDecoder(r.Body).Decode(&sw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if sw.OldUserID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("'old_user_id' can't be lower or equal zero, got %d", sw.OldUserID))

		return
	}

	if sw.NewUserID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("'new_user_id' can't be lower or equal zero, got %d", sw.NewUserID))

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
		WriteError(w, fmt.Errorf("switch: %w", err))

		return
	}

	w.Write(success)
}
