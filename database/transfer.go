package database

import (
	"fmt"
	"time"
)

func (sql DB) Transfer(fromID, toID, amount int64, description string) error {
	if amount < 0 {
		return fmt.Errorf("amount can't be lower than zero, got %d", amount)
	}

	sender, senderID, err := sql.GetBalance(fromID)
	if err != nil {
		return fmt.Errorf("get sender balance: %w", err)
	}

	if amount > sender {
		return fmt.Errorf("insufficient funds: missing %.2f", float64(amount-sender)/100)
	}

	var receiver, receiverID int64

	rows, err := sql.conn.Query("SELECT balance,balance_id FROM balances WHERE user_id = ?", toID)
	if err != nil {
		return fmt.Errorf("get balance: query: %w", err)
	}
	if rows.Next() { // If balance is found
		err = rows.Scan(&receiver, &receiverID)
		if err != nil {
			return fmt.Errorf("get balance: scan: %w", err)
		}
	} else { // Create new account
		r, err := sql.conn.Exec("INSERT INTO balances (user_id,balance) VALUES (?,?)", toID, amount)
		if err != nil {
			return fmt.Errorf("new balance: insert balance (user_id %d): %w", toID, err)
		}

		receiverID, err = r.LastInsertId()
		if err != nil {
			return fmt.Errorf("new balance: get LastInsertId: %w", err)
		}
		receiver = -1
	}

	tx, err := sql.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transactions: %w", err)
	}

	sender -= amount
	_, err = tx.Exec("UPDATE balances SET balance = ? WHERE balance_id = ?", sender, senderID)
	if err != nil {
		return fmt.Errorf(
			"update sender balance (user_id %d, change %d, new balance %d): %w",
			fromID, -amount, sender, err,
		)
	}

	if receiver >= 0 {
		receiver += amount
		_, err = tx.Exec("UPDATE balances SET balance = ? WHERE balance_id = ?", receiver, receiverID)
		if err != nil {
			return fmt.Errorf(
				"update receiver balance (user_id %d, change %d, new balance %d): %w",
				toID, amount, sender, err,
			)
		}
	}

	_, err = tx.Exec(`INSERT INTO transactions (balance_id,from_id,action,description,date)
	VALUES (?,?,?,?,?)`, receiverID, senderID, amount, description, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
