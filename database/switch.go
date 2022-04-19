package database

import (
	"context"
	"fmt"
)

func (sql DB) Switch(oldUserID, newUserID int64) error {
	ctx := context.Background()

	tx, err := sql.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transactions: %w", err)
	}
	defer tx.Rollback(ctx)

	_, balanceID, err := getBalanceForUpdate(tx, oldUserID)
	if err != nil {
		return fmt.Errorf("GetBalance(old %d): %w", oldUserID, err)
	}

	_, _, err = sql.GetBalance(ctx, newUserID)
	if err == nil {
		return fmt.Errorf("balance with user_id %d already exists", newUserID)
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET user_id = $1 WHERE balance_id = $2", newUserID, balanceID)
	if err != nil {
		return fmt.Errorf("switch user id (old %d - new %d): %w", oldUserID, newUserID, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
