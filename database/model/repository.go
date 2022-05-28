package model

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	GetBalance(ctx context.Context, userID int64) (balance, balanceID int64, err error)
	GetBalanceForUpdate(ctx context.Context, db pgx.Tx, userID int64) (balance, balanceID int64, err error)
	ChangeBalance(ctx context.Context, userID, change int64, description string) error

	Transfer(ctx context.Context, fromID, toID, amount int64, description string) error
	GetTransfers(ctx context.Context, userID, offset, limit int64, sort string) ([]Transaction, error)

	Switch(ctx context.Context, oldUserID, newUserID int64) error
	Delete(ctx context.Context, userID int64) error
}

// Transaction
// @Description User Transaction
type Transaction struct {
	TransactionID int64            `json:"transaction_id"`
	BalanceID     int64            `json:"balance_id"`
	FromID        string           `json:"from_id"`
	Action        int64            `json:"action"`
	Date          pgtype.Timestamp `json:"date,string" swaggertype:"string"`
	Description   string           `json:"description"`
}
