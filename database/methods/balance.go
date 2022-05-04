package methods

import (
	"autumn-2021-intern-assignment/public"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Methods struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) Methods {
	return Methods{conn}
}

func (sql Methods) GetBalance(ctx context.Context, userID int64) (balance, balanceID int64, err error) {

	err = sql.conn.QueryRow(ctx, "SELECT balance,balance_id FROM balances WHERE user_id = $1", userID).
		Scan(&balance, &balanceID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, -1, nil
		}

		return -1, -1, public.NewInternal(fmt.Errorf("query: %w", err))
	}

	return
}

func (sql Methods) ChangeBalance(ctx context.Context, userID, change int64, description string) error {

	t, err := sql.conn.Begin(ctx)
	if err != nil {
		return public.NewInternal(fmt.Errorf("begin transaction: %w", err))
	}
	defer t.Rollback(ctx)

	balance, balanceID, err := sql.GetBalanceForUpdate(ctx, t, userID)

	if err != nil {
		return err
	}

	if balanceID < 0 {
		if change < 0 {
			return fmt.Errorf("change (%d) is below zero, balance creating is forbidden", change)
		}

		err = t.QueryRow(ctx,
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

	balance += change

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

func (Methods) GetBalanceForUpdate(ctx context.Context, tx pgx.Tx, userID int64) (balance, balanceID int64, err error) {
	err = tx.QueryRow(ctx, "SELECT balance,balance_id FROM balances WHERE user_id = $1", userID).
		Scan(&balance, &balanceID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, -1, nil
		}

		return -1, -1, public.NewInternal(fmt.Errorf("get balance: query: %w", err))
	}

	return
}
