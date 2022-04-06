package methods

import (
	"encoding/json"

	"autumn-2021-intern-assignment/database"
	"github.com/valyala/fasthttp"
)

type transferJSON struct {
	ToID        int64  `json:"to_id"`
	FromID      int64  `json:"from_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

func Transfer(db *database.DB, ctx *fasthttp.RequestCtx) {
	var trans transferJSON

	err := json.Unmarshal(ctx.PostBody(), &trans)
	if err != nil {
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	if trans.ToID <= 0 {
		ctx.Write(jsonError("wrong 'to_id' field value: %d", trans.ToID))
		return
	}
	if trans.FromID <= 0 {
		ctx.Write(jsonError("wrong 'from_id' field value: %d", trans.FromID))
		return
	}
	if trans.Amount <= 0 {
		ctx.Write(jsonError("wrong 'amount' field value: can't be lower or equal zero, got %d", trans.Amount))
		return
	}

	err = db.Transfer(trans.FromID, trans.ToID, trans.Amount, trans.Description)
	if err != nil {
		ctx.Write(jsonError("transfer: %w", err))
		return
	}

	ctx.Write(jsonSuccess())
}
