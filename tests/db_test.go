package server_test

import (
	"autumn-2021-intern-assignment/database"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

func FuzzDatabase(f *testing.F) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		f.Fatal(err)
	}

	db := database.NewFromConnect(conn)

	f.Add(1000, 1, 2)
	f.Fuzz(ChangeTest(db, mock))

	defer conn.Close()
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func ChangeTest(db *database.DB, mock sqlmock.Sqlmock) any {
	return func(t *testing.T, balance int, userID1, userID2 int) {
		userID1, userID2 = Abs(userID1), Abs(userID2)

		// Get Balance 1
		mock.ExpectQuery("SELECT balance,balance_id FROM balances WHERE user_id =").
			WithArgs(userID1).
			WillReturnRows(&sqlmock.Rows{})

		db.GetBalance(int64(userID1))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("GetBalance (userID1 %d): %s", userID1, err)
		}

		// Get Balance 2
		mock.ExpectQuery("SELECT balance,balance_id FROM balances WHERE user_id =").
			WithArgs(userID2).
			WillReturnRows(&sqlmock.Rows{})

		db.GetBalance(int64(userID2))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("GetBalance (userID2 %d): %s", userID2, err)
		}

		// Change balance with create
		mock.ExpectQuery("SELECT balance,balance_id FROM balances WHERE user_id =").
			WithArgs(userID1).
			WillReturnRows(&sqlmock.Rows{})

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO balances (.+) VALUES (.+)").
			WithArgs(userID1, balance).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO transactions (.+) VALUES (.+)").
			WithArgs(1, int64(balance), "test", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := db.ChangeBalance(int64(userID1), int64(balance), "test")
		if err != nil {
			t.Error(err)
		}

		// Get Balance 1
		mock.ExpectQuery("SELECT balance,balance_id FROM balances WHERE user_id =").
			WithArgs(userID1).
			WillReturnRows(sqlmock.NewRows([]string{"balance", "balance_id"}).AddRow(balance, 1))
		{
			bal, bID, err := db.GetBalance(int64(userID1))
			if err != nil {
				t.Error(err)
			}
			if bID != 1 {
				t.Errorf("GetBalance after changing(user_id %d): wrong values: got balance_id %d instead of %d",
					userID1, bID, 1,
				)
			}
			if bal != int64(balance) {
				t.Errorf("GetBalance after changing(user_id %d): wrong values: got balance %d instead of %d",
					userID1, bal, balance,
				)
			}
		}
	}
}
