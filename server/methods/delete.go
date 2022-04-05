package methods

import (
	"autumn-2021-intern-assignment/database"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type deleteJSON struct {
	UserID    int64 `json:"user_id"`
	BalanceID int64 `json:"balance_id"`
}

func Delete(db *database.DB, ctx *fasthttp.RequestCtx) {
	var del deleteJSON

	err := json.Unmarshal(ctx.PostBody(), &del)
	if err != nil {
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	switch {
	case del.BalanceID >= 1:
		del.UserID = 0

	case del.UserID >= 1:
		del.BalanceID = 0
	default:
		ctx.Write(
			jsonError("both 'user_id' and 'balance_id' fields cant be wrong format, got %d %d", del.BalanceID, del.UserID),
		)
		return
	}

	err = db.Delete(del.BalanceID, del.UserID)
	if err != nil {
		ctx.Write(jsonError("delete balance (user_id %d, balance_id %d): %w", del.UserID, del.BalanceID, err))
		return
	}

	ctx.Write(jsonSuccess())
}
