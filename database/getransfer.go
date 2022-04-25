package database

import (
	"context"
	types "database/sql"
	"fmt"
	"strconv"

	"autumn-2021-intern-assignment/public"
	"github.com/jackc/pgtype"
)

var sorts = map[string]string{
	"DATE_DESC": "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3",
	"DATE_ASC":  "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY date ASC LIMIT $2 OFFSET $3",
	"SUM_DESC":  "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY action DESC LIMIT $2 OFFSET $3",
	"SUM_ASC":   "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY action ASC LIMIT $2 OFFSET $3",
}

type Transaction struct {
	TransactionID int64            `json:"transaction_id"`
	BalanceID     int64            `json:"balance_id"`
	FromID        string           `json:"from_id"`
	Action        int64            `json:"action"`
	Date          pgtype.Timestamp `json:"date"`
	Description   string           `json:"description"`
}

func (sql DB) GetTransfers(userID, offset, limit int64, sort string) ([]Transaction, error) {
	ctx := context.Background()

	_, balanceID, err := sql.GetBalance(ctx, userID)
	if err != nil {
		return nil, public.NewInternal(fmt.Errorf("get balance (id %d): %w", userID, err))
	}

	command, ok := sorts[sort]

	if !ok {
		return nil, fmt.Errorf("sort %s not supported", sort)
	}

	rows, err := sql.conn.Query(ctx, command, balanceID, limit, offset)
	if err != nil {
		return nil, public.NewInternal(fmt.Errorf("select query: %w", err))
	}
	defer rows.Close()

	var (
		trs  []Transaction
		t    Transaction
		from types.NullInt64
	)

	for rows.Next() {
		err = rows.Scan(&t.TransactionID, &t.BalanceID, &from, &t.Action, &t.Date, &t.Description)

		if err != nil {
			return nil, public.NewInternal(fmt.Errorf("scan transaction: %w", err))
		}

		if from.Valid {
			t.FromID = strconv.FormatInt(from.Int64, 10)
		} else {
			t.FromID = ""
		}

		trs = append(trs, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, public.NewInternal(fmt.Errorf("rows: %w", err))
	}

	return trs, nil
}
