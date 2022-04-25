package methods

import (
	"encoding/json"
	"errors"
	"net/http"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/public"
	"github.com/valyala/fasthttp"
)

type viewJSON struct {
	User   int64  `json:"user_id"`
	Sort   string `json:"sort"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

type viewRetJSON struct {
	Ok           bool                   `json:"ok"`
	Transactions []database.Transaction `json:"transactions"`
}

func View(db *database.DB, ctx *fasthttp.RequestCtx) {
	var view viewJSON

	err := json.Unmarshal(ctx.PostBody(), &view)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "decoding json: %w", err)

		return
	}

	if view.User <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "wrong 'user' field value: %d", view.User)

		return
	}

	if view.Limit < 1 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "wrong 'limit' field value: cant be lower than 1, got %d", view.Limit)

		return
	}

	if view.Offset < 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonError(ctx, "wrong 'offset' field value: cant be lower than 0 got %d", view.Offset)

		return
	}

	if view.Sort == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		jsonErrorString(ctx, "'sort' field value cant be empty")

		return
	}

	trans, err := db.GetTransfers(view.User, view.Offset, view.Limit, view.Sort)
	if err != nil {
		if errors.As(err, public.ErrInternal) {
			ctx.SetStatusCode(http.StatusInternalServerError)
		} else {
			ctx.SetStatusCode(http.StatusUnprocessableEntity)
		}
		jsonError(ctx, "get transfers: %w", err)

		return
	}

	data, err := json.Marshal(viewRetJSON{
		Ok:           true,
		Transactions: trans,
	})
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		jsonError(ctx, "encoding json: %w", err)

		return
	}

	ctx.Write(data)
}
