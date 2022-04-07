package methods

import (
	"autumn-2021-intern-assignment/database"
	"encoding/json"
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
		ctx.Write(jsonError("decoding json: %w", err))
		return
	}

	if sw.OldUserID <= 0 {
		ctx.Write(jsonError("'old_user_id' can't be lower or equal zero, got %d", sw.OldUserID))
		return
	}

	if sw.NewUserID <= 0 {
		ctx.Write(jsonError("'new_user_id' can't be lower or equal zero, got %d", sw.NewUserID))
		return
	}

	err = db.Switch(sw.OldUserID, sw.NewUserID)
	if err != nil {
		ctx.Write(jsonError("db.Switch(old %d - new %d): %w", sw.OldUserID, sw.NewUserID, err))
		return
	}

	ctx.Write(jsonSuccess())
}
