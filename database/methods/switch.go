package methods

import (
	"context"
	"fmt"

	"autumn-2021-intern-assignment/public"
)

func (sql Methods) Switch(ctx context.Context, oldUserID, newUserID int64) error {

	tx, err := sql.conn.Begin(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("begin transactions: %w", err))
	}
	defer tx.Rollback(ctx)

	_, balanceID, err := sql.GetBalanceForUpdate(ctx, tx, oldUserID)
	if err != nil {
		return public.NewInternal(fmt.Errorf("GetBalance(old %d): %w", oldUserID, err))
	}

	_, bufID, err := sql.GetBalance(ctx, newUserID)
	if err != nil {
		return err
	}
	if bufID > 0 {
		return fmt.Errorf("balance with user_id %d already exists", newUserID)
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET user_id = $1 WHERE balance_id = $2", newUserID, balanceID)
	if err != nil {
		return public.NewInternal(fmt.Errorf("switch user id (old %d - new %d): %w", oldUserID, newUserID, err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}
