package methods

import (
	"autumn-2021-intern-assignment/public"
	"context"
	"fmt"
)

func (m Methods) Delete(ctx context.Context, userID int64) error {
	tag, err := m.conn.Exec(ctx, "DELETE FROM balances WHERE user_id = $1", userID)
	if err != nil {
		return public.NewInternal(err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("balance with user id %d not found", userID)
	}

	return nil
}
