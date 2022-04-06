package methods

import (
	"encoding/json"

	"autumn-2021-intern-assignment/database"
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
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	if change.User <= 0 {
		ctx.Write(jsonError("wrong 'user' field value: %d", change.User))
		return
	}

	if change.Change == 0 {
		ctx.Write(jsonErrorString("wrong 'change' field value: can't be zero"))
		return
	}

	if change.Description == "" {
		ctx.Write(jsonError("'description' field value can't be empty"))
		return
	}

	err = db.ChangeBalance(change.User, change.Change, change.Description)
	if err != nil {
		ctx.Write(jsonError("change balance: %w", err))
		return
	}

	ctx.Write(jsonSuccess())
}
