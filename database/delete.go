package database

import "fmt"

func (sql DB) Delete(balanceID, userID int64) error {
	var err error

	if balanceID < 1 {
		_, balanceID, err = sql.GetBalance(userID)
		if err != nil {
			return fmt.Errorf("get balance (userID %d): %w", userID, err)
		}
	}

	t, err := sql.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = t.Exec("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return fmt.Errorf("disable foreign key checks: %w", err)
	}

	_, err = t.Exec("DELETE FROM balances WHERE balance_id = ?", balanceID)
	if err != nil {
		return fmt.Errorf("deleting balance (id %d): %w", balanceID, err)
	}

	_, err = t.Exec("SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return fmt.Errorf("enable foreign key checks: %w", err)
	}

	err = t.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
