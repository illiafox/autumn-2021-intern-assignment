package methods

import (
	"context"
	"fmt"
	"time"

	"autumn-2021-intern-assignment/public"
)

func (sql Methods) Transfer(ctx context.Context, fromID, toID, amount int64, description string) error {
	if amount < 0 {
		return fmt.Errorf("amount can't be lower than zero, got %d", amount)
	}

	tx, err := sql.conn.Begin(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("begin transactions: %w", err))
	}
	defer tx.Rollback(ctx)

	sender, senderID, err := sql.GetBalanceForUpdate(ctx, tx, fromID)
	if err != nil {
		return public.NewInternal(fmt.Errorf("get sender balance: %w", err))
	}

	if amount > sender {
		return fmt.Errorf("insufficient funds: missing %.2f", float64(amount-sender)/100)
	}

	_, receiverID, err := sql.GetBalanceForUpdate(ctx, tx, toID)
	if err != nil {
		return err
	}

	if receiverID < 0 {
		err = tx.QueryRow(
			ctx,
			"INSERT INTO balances (user_id,balance) VALUES ($1,$2) RETURNING balance_id", toID, amount,
		).Scan(&receiverID)

		if err != nil {
			return public.NewInternal(fmt.Errorf("new balance: insert balance (user_id %d): %w", toID, err))
		}
	} else {
		_, err = tx.Exec(ctx, "UPDATE balances SET balance = balance + amount WHERE balance_id = $2", amount, receiverID)

		if err != nil {
			return public.NewInternal(fmt.Errorf(
				"update receiver balance (user_id %d, change %d, new balance %d): %w",
				toID, amount, sender, err,
			))
		}
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = balance - $1 WHERE balance_id = $2", amount, senderID)

	if err != nil {
		return public.NewInternal(fmt.Errorf(
			"update sender balance (user_id %d, change %d, new balance %d): %w",
			fromID, -amount, sender, err,
		))
	}

	_, err = tx.Exec(ctx, `INSERT INTO transactions (balance_id,from_id,action,description,date)
	VALUES ($1,$2,$3,$4,$5)`, receiverID, senderID, amount, description, time.Now().UTC())
	if err != nil {
		return public.NewInternal(fmt.Errorf("insert transaction: %w", err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}
