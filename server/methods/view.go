package methods

import (
	"autumn-2021-intern-assignment/database"
	"encoding/json"
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
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	if view.User <= 0 {
		ctx.Write(jsonError("wrong 'user' field value: %d", view.User))
		return
	}
	if view.Limit < 1 {
		ctx.Write(jsonError("wrong 'limit' field value: cant be lower than 1, got %d", view.User))
		return
	}
	if view.Offset < 0 {
		ctx.Write(jsonError("wrong 'offset' field value: cant be lower than 0 got %d", view.User))
		return
	}
	if view.Sort == "" {
		ctx.Write(jsonErrorString("'sort' field value cant be empty"))
		return
	}

	trans, err := db.GetTransfers(view.User, view.Offset, view.Limit, view.Sort)
	if err != nil {
		ctx.Write(jsonError("get transfers: %w", err))
		return
	}

	data, err := json.Marshal(viewRetJSON{
		Ok:           true,
		Transactions: trans,
	})
	if err != nil {
		ctx.Write(jsonError("encoding json: %w", err))
		return
	}

	ctx.Write(data)
}
