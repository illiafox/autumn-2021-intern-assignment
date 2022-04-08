package database

import "fmt"

func (sql DB) Switch(oldUserID, newUserID int64) error {

	tx, err := sql.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transactions: %w", err)
	}
	defer tx.Rollback()

	_, balanceID, err := GetBalanceForUpdate(tx, oldUserID)
	if err != nil {
		return fmt.Errorf("GetBalance(old %d): %w", oldUserID, err)
	}

	_, _, err = sql.GetBalance(newUserID)
	if err == nil {
		return fmt.Errorf("balance with user_id %d already exists", newUserID)
	}

	_, err = tx.Exec("UPDATE balances SET user_id = $1 WHERE balance_id = $2", newUserID, balanceID)
	if err != nil {
		return fmt.Errorf("switch user id (old %d - new %d): %w", oldUserID, newUserID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
