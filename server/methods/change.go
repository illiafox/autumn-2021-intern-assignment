package methods

import (
	"encoding/json"
	"errors"
	"net/http"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/public"
	"github.com/valyala/fasthttp"
)

type changeJSON struct {
	User        int64  `json:"user_id"`
	Change      int64  `json:"change"`
	Description string `json:"description"`
}

func Change(db *database.DB, ctx *fasthttp.RequestCtx) {
	var change changeJSON

	err := json.Unmarshal(ctx.PostBody(), &change)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "decoding json: %w", err)

		return
	}

	if change.User <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "wrong 'user' field value: %d", change.User)

		return
	}

	if change.Change == 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonErrorString(ctx, "wrong 'change' field value: can't be zero")

		return
	}

	if change.Description == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "'description' field value can't be empty")

		return
	}

	err = db.ChangeBalance(change.User, change.Change, change.Description)
	if err != nil {
		if errors.As(err, public.ErrInternal) {
			ctx.SetStatusCode(http.StatusUnprocessableEntity)
		} else {
			ctx.SetStatusCode(http.StatusNotAcceptable)
		}
		jsonError(ctx, "change balance: %w", err)

		return
	}

	jsonSuccess(ctx)
}
