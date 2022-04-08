package database

import (
	sqlpack "database/sql"
	"fmt"
	"strconv"
)

var sorts = map[string]string{
	"DATE_DESC": "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3",
	"DATE_ASC":  "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY date ASC LIMIT $2 OFFSET $3",
	"SUM_DESC":  "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY action DESC LIMIT $2 OFFSET $3",
	"SUM_ASC":   "SELECT * FROM transactions WHERE balance_id = $1 ORDER BY action ASC LIMIT $2 OFFSET $3",
}

type Transaction struct {
	TransactionID int64  `json:"transaction_id"`
	BalanceID     int64  `json:"balance_id"`
	FromID        string `json:"from_id"`
	Action        int64  `json:"action"`
	Date          string `json:"date"`
	Description   string `json:"description"`
}

func (sql DB) GetTransfers(userID, offset, limit int64, sort string) ([]Transaction, error) {
	_, balanceID, err := sql.GetBalance(userID)
	if err != nil {
		return nil, fmt.Errorf("get balance (id %d): %w", userID, err)
	}

	command, ok := sorts[sort]

	if !ok {
		return nil, fmt.Errorf("sort %s not supported", sort)
	}

	rows, err := sql.conn.Query(command, balanceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select query: %w", err)
	}
	defer rows.Close()

	var (
		trs     []Transaction
		t       Transaction
		fromBuf sqlpack.NullInt64
	)

	for rows.Next() {
		err = rows.Scan(&t.TransactionID, &t.BalanceID, &fromBuf, &t.Action, &t.Date, &t.Description)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		if fromBuf.Valid {
			t.FromID = strconv.FormatInt(fromBuf.Int64, 10)
		} else {
			t.FromID = ""
		}

		trs = append(trs, t)
	}

	return trs, nil
}
