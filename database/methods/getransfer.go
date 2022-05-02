package methods

import (
	"context"
	types "database/sql"
	"fmt"
	"strconv"
	"strings"

	"autumn-2021-intern-assignment/database/model"
	"autumn-2021-intern-assignment/public"
)

var (
	sorts = map[string]string{
		"DATE_DESC": "date DESC",
		"DATE_ASC":  "date ASC",
		"SUM_DESC":  "action DESC",
		"SUM_ASC":   "action ASC",
	}

	available = "sort '%s' not supported, available: " + func() string {
		keys := make([]string, 0, len(sorts))
		for k := range sorts {
			keys = append(keys, k)
		}

		return strings.Join(keys, " ")
	}()
)

func (sql Methods) GetTransfers(ctx context.Context, userID, offset, limit int64,
	sort string) ([]model.Transaction, error) {

	order, ok := sorts[sort]
	if !ok {
		return nil, fmt.Errorf(available, sort)
	}

	_, balanceID, err := sql.GetBalance(ctx, userID)
	if err != nil {
		return nil, public.NewInternal(fmt.Errorf("get balance (id %d): %w", userID, err))
	}

	rows, err := sql.conn.Query(ctx,
		"SELECT * FROM transactions WHERE balance_id = $1 ORDER BY "+order+" LIMIT $2 OFFSET $3",
		balanceID, limit, offset)

	if err != nil {
		return nil, public.NewInternal(fmt.Errorf("select query: %w", err))
	}
	defer rows.Close()

	var (
		trs  []model.Transaction
		t    model.Transaction
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
