package database

import (
	"fmt"
	"time"
)

func (sql DB) Transfer(fromID, toID, amount int64, description string) error {
	if amount < 0 {
		return fmt.Errorf("amount can't be lower than zero, got %d", amount)
	}

	tx, err := sql.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transactions: %w", err)
	}
	defer tx.Rollback()

	sender, senderID, err := GetBalanceForUpdate(tx, fromID)
	if err != nil {
		return fmt.Errorf("get sender balance: %w", err)
	}

	if amount > sender {
		return fmt.Errorf("insufficient funds: missing %.2f", float64(amount-sender)/100)
	}

	var receiver, receiverID int64

	rows, err := tx.Query("SELECT balance,balance_id FROM balances WHERE user_id = $1 FOR UPDATE", toID)
	if err != nil {
		return fmt.Errorf("get balance: query: %w", err)
	}
	defer rows.Close()

	if rows.Next() { // If balance is found
		err = rows.Scan(&receiver, &receiverID)
		if err != nil {
			return fmt.Errorf("get balance: scan: %w", err)
		}

	} else { // Create new account
		err = tx.QueryRow(
			"INSERT INTO balances (user_id,balance) VALUES ($1,$2) RETURNING balance_id", toID, amount,
		).Scan(&receiverID)

		if err != nil {
			return fmt.Errorf("new balance: insert balance (user_id %d): %w", toID, err)
		}

		receiver = -1
	}

	sender -= amount
	_, err = tx.Exec("UPDATE balances SET balance = $1 WHERE balance_id = $2", sender, senderID)

	if err != nil {
		return fmt.Errorf(
			"update sender balance (user_id %d, change %d, new balance %d): %w",
			fromID, -amount, sender, err,
		)
	}

	if receiver >= 0 {
		receiver += amount
		_, err = tx.Exec("UPDATE balances SET balance = $1 WHERE balance_id = $2", receiver, receiverID)

		if err != nil {
			return fmt.Errorf(
				"update receiver balance (user_id %d, change %d, new balance %d): %w",
				toID, amount, sender, err,
			)
		}
	}

	_, err = tx.Exec(`INSERT INTO transactions (balance_id,from_id,action,description,date)
	VALUES ($1,$2,$3,$4,$5)`, receiverID, senderID, amount, description, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
