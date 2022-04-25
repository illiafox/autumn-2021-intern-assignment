package methods

import (
	"encoding/json"
	"errors"
	"net/http"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/public"
	"github.com/valyala/fasthttp"
)

type switchJSON struct {
	OldUserID int64 `json:"old_user_id"`
	NewUserID int64 `json:"new_user_id"`
}

func Switch(db *database.DB, ctx *fasthttp.RequestCtx) {
	var sw switchJSON

	err := json.Unmarshal(ctx.PostBody(), &sw)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "decoding json: %w", err)

		return
	}

	if sw.OldUserID <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "'old_user_id' can't be lower or equal zero, got %d", sw.OldUserID)

		return
	}

	if sw.NewUserID <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "'new_user_id' can't be lower or equal zero, got %d", sw.NewUserID)

		return
	}

	err = db.Switch(sw.OldUserID, sw.NewUserID)
	if err != nil {
		if errors.As(err, public.ErrInternal) {
			ctx.SetStatusCode(http.StatusUnprocessableEntity)
		} else {
			ctx.SetStatusCode(http.StatusNotAcceptable)
		}
		jsonError(ctx, "switch: %w", err)

		return
	}

	jsonSuccess(ctx)
}
