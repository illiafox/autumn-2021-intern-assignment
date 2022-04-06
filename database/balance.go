package database

import (
	"fmt"
	"time"
)

func (sql DB) GetBalance(userID int64) (balance, balanceID int64, err error) {
	rows, err := sql.conn.Query("SELECT balance,balance_id FROM balances WHERE user_id = ?", userID)
	if err != nil {
		return -1, -1, fmt.Errorf("query: %w", err)
	}

	if rows.Next() { // If balance is found
		err = rows.Scan(&balance, &balanceID)
		if err != nil {
			return -1, -1, fmt.Errorf("scan: %w", err)
		}

		return
	}

	return -1, -1, fmt.Errorf("balance with user id %d not found", userID)
}

func (sql DB) ChangeBalance(userID, change int64, description string) error {
	rows, err := sql.conn.Query("SELECT balance,balance_id FROM balances WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("get balance: query: %w", err)
	}

	t, err := sql.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	var balance, balanceID int64

	if rows.Next() { // If balance is found
		err = rows.Scan(&balance, &balanceID)
		if err != nil {
			return fmt.Errorf("get balance: scan: %w", err)
		}

		balance += change
	} else { // Create new account
		if change < 0 {
			return fmt.Errorf("change (%d) is below zero, balance creating is forbidden", change)
		}

		r, err := sql.conn.Exec("INSERT INTO balances (user_id,balance) VALUES (?,?)", userID, change)
		if err != nil {
			return fmt.Errorf(
				"new balance: insert balance (user_id %d, change %d): %w",
				userID, change, err,
			)
		}

		balanceID, err = r.LastInsertId()
		if err != nil {
			return fmt.Errorf("new balance: get LastInsertId: %w", err)
		}
		goto trans
	}

	if balance < 0 {
		return fmt.Errorf("insufficient funds: missing %.2f", float64(-balance)/100)
	}

	_, err = t.Exec("UPDATE balances SET balance = ? WHERE balance_id = ?", balance, balanceID)
	if err != nil {
		return fmt.Errorf(
			"update balance (id %d, change %d, new balance %d): %w",
			balanceID, change, balance, err,
		)
	}
trans:
	_, err = t.Exec(`INSERT INTO transactions (balance_id,action,description,date) 
	VALUES (?,?,?,?)`, balanceID, change, description, time.Now().Format(time.RFC3339))

	if err != nil {
		return fmt.Errorf(
			"insert transaction (user_id %d, change %d): %w",
			userID, change, err,
		)
	}

	err = t.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
