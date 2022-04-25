package database

import (
	"context"
	"fmt"
	"time"

	"autumn-2021-intern-assignment/public"
	"github.com/jackc/pgx/v4"
)

func (sql DB) GetBalance(ctx context.Context, userID int64) (balance, balanceID int64, err error) {

	rows, err := sql.conn.Query(ctx, "SELECT balance,balance_id FROM balances WHERE user_id = $1", userID)

	if err != nil {
		return -1, -1, public.NewInternal(fmt.Errorf("query: %w", err))
	}
	defer rows.Close()

	if rows.Next() { // If balance is found
		err = rows.Scan(&balance, &balanceID)
		if err != nil {
			return -1, -1, public.NewInternal(fmt.Errorf("scan: %w", err))
		}

		return
	}

	err = rows.Err()
	if err != nil {
		return -1, -1, public.NewInternal(fmt.Errorf("rows: %w", err))
	}

	return -1, -1, fmt.Errorf("balance with user id %d not found", userID)
}

func (sql DB) ChangeBalance(userID, change int64, description string) error {
	ctx := context.Background()

	t, err := sql.conn.Begin(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("begin transaction: %w", err))
	}
	defer t.Rollback(ctx)

	rows, err := t.Query(ctx, "SELECT balance,balance_id FROM balances WHERE user_id = $1 FOR UPDATE", userID)
	if err != nil {
		return public.NewInternal(fmt.Errorf("get balance: query: %w", err))
	}
	defer rows.Close()
	var balance, balanceID int64

	if rows.Next() { // If balance is found
		err = rows.Scan(&balance, &balanceID)
		if err != nil {
			return public.NewInternal(fmt.Errorf("get balance: scan: %w", err))
		}

		balance += change
	} else { // Create new account
		err = rows.Err()
		if err != nil {
			return public.NewInternal(fmt.Errorf("rows: %w", err))
		}

		if change < 0 {
			return fmt.Errorf("change (%d) is below zero, balance creating is forbidden", change)
		}

		err = t.QueryRow(
			ctx,
			"INSERT INTO balances (user_id,balance) VALUES ($1,$2) RETURNING balance_id", userID, change,
		).Scan(&balanceID)

		if err != nil {
			return public.NewInternal(fmt.Errorf(
				"new balance: insert balance (user_id %d, change %d)): %w",
				userID, change, err,
			))
		}

		goto final
	}

	if balance < 0 {
		return fmt.Errorf("insufficient funds: missing %.2f", float64(-balance)/100)
	}

	_, err = t.Exec(ctx, "UPDATE balances SET balance = $1 WHERE balance_id = $2", balance, balanceID)
	if err != nil {
		return public.NewInternal(fmt.Errorf(
			"update balance (id %d, change %d, new balance %d): %w",
			balanceID, change, balance, err,
		))
	}
final:
	_, err = t.Exec(ctx, `INSERT INTO transactions (balance_id,action,description,date) 
	VALUES ($1,$2,$3,$4)`, balanceID, change, description, time.Now().UTC())

	if err != nil {
		return public.NewInternal(fmt.Errorf(
			"insert transaction (user_id %d, change %d): %w",
			userID, change, err,
		))
	}

	err = t.Commit(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}

func getBalanceForUpdate(db pgx.Tx, userID int64) (balance, balanceID int64, err error) {
	rows, err := db.Query(context.Background(), "SELECT balance,balance_id FROM balances WHERE user_id = $1", userID)

	if err != nil {
		return -1, -1, public.NewInternal(fmt.Errorf("query: %w", err))
	}
	defer rows.Close()

	if rows.Next() { // If balance is found
		err = rows.Scan(&balance, &balanceID)
		if err != nil {
			return -1, -1, public.NewInternal(fmt.Errorf("scan: %w", err))
		}

		return
	}

	err = rows.Err()
	if err != nil {
		return -1, -1, public.NewInternal(fmt.Errorf("rows: %w", err))
	}

	return -1, -1, fmt.Errorf("balance with user id %d not found", userID)
}
