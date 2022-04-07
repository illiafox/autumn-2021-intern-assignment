package database

import "fmt"

func (sql DB) Switch(oldUserID, newUserID int64) error {

	_, balanceID, err := sql.GetBalance(oldUserID)
	if err != nil {
		return fmt.Errorf("GetBalance(old %d): %w", oldUserID, err)
	}

	_, _, err = sql.GetBalance(newUserID)
	if err == nil {
		return fmt.Errorf("balance with user_id %d already exists", newUserID)
	}

	_, err = sql.conn.Exec("UPDATE balances SET user_id = ? WHERE balance_id = ?", newUserID, balanceID)
	if err != nil {
		return fmt.Errorf("switch user id (old %d - new %d): %w", oldUserID, newUserID, err)
	}

	return nil
}
