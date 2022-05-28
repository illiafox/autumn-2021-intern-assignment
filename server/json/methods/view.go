package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"autumn-2021-intern-assignment/database/model"
	"autumn-2021-intern-assignment/public"
)

// View
// @Description User ID and view options
type View struct {
	User   int64  `json:"user_id"`
	Sort   string `json:"sort"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

type ViewOut struct {
	Ok           bool                `json:"ok"`
	Transactions []model.Transaction `json:"transactions"`
}

// View godoc
// @Summary      View transactions
// @Description  View user transactions
// @Accept       json
// @Produce      json
// @Param        input body 	View true "User ID and view options"
// @Success      200  {object}  ViewOut
// @Failure      400  {object}  Error
// @Failure      422  {object}  Error
// @Failure      500  {object}  Error
// @Router       /view [get]
func (m Methods) View(w http.ResponseWriter, r *http.Request) {
	var view View

	err := json.NewDecoder(r.Body).Decode(&view)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("decoding json: %w", err))

		return
	}

	if view.User <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'user' field value: %d", view.User))

		return
	}

	if view.Limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'limit' field value: cant be lower than 1, got %d", view.Limit))

		return
	}

	if view.Offset < 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, fmt.Errorf("wrong 'offset' field value: cant be lower than 0 got %d", view.Offset))

		return
	}

	if view.Sort == "" {
		w.WriteHeader(http.StatusBadRequest)
		WriteString(w, "'sort' field value cant be empty")

		return
	}

	ctx := context.Background()

	trans, err := m.db.GetTransfers(ctx, view.User, view.Offset, view.Limit, view.Sort)
	if err != nil {
		if public.AsInternal(err) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		WriteError(w, fmt.Errorf("get transfers: %w", err))

		return
	}

	json.NewEncoder(w).Encode(ViewOut{true, trans})
}
