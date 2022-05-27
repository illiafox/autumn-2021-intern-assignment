package db

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/utils/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite

	db       *database.Database
	truncate func(context.Context) error
}

func (suite *DBSuite) SetupSuite() {
	rand.Seed(time.Now().UnixNano())

	conf := &config.Postgres{
		User:     "server",
		Pass:     "M5F3wWtFxkQ8Ra4n",
		DbName:   "avito_test",
		IP:       "127.0.0.1",
		Port:     "5432",
		Protocol: "tcp",
	}

	r := suite.Require()

	// //

	err := config.ReadEnv(conf)
	r.NoError(err, "read environment")

	// //

	ctx := context.Background()

	// //

	pool, err := pgxpool.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable",
			conf.User,
			conf.Pass,
			conf.IP,
			conf.Port,
			conf.DbName,
		),
	)
	r.NoError(err, "connect to database")

	// //

	err = pool.Ping(ctx)
	r.NoError(err, "ping connection")

	count := 0

	// //

	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM balances").Scan(&count)
	r.NoError(err, "select count from balances")

	r.Equal(0, count, "'balances' table is not empty")

	// //

	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	r.NoError(err, "select count from balances")

	r.Equal(0, count, "'transactions' table is not empty")

	// //

	suite.db = database.NewDatabase(pool)
	suite.truncate = func(ctx context.Context) error {
		_, err := pool.Exec(ctx, "TRUNCATE TABLE transactions,balances")

		return err
	}
}

func (suite *DBSuite) TearDownSuite() {
	if suite.db == nil {
		return
	}

	err := suite.truncate(context.Background())
	suite.Assert().NoError(err, "truncate tables")

	suite.db.Close()
}

func (suite *DBSuite) TestChange() {
	ctx := context.Background()

	id, change := rand.Int63n(100000), rand.Int63n(100000)
	a := suite.Assert()

	err := suite.db.ChangeBalance(ctx, id, change, "deposit")
	a.NoError(err, "change balance with creation (user_id %d, change %d)", id, change)

	balance, _, err := suite.db.GetBalance(ctx, id)
	a.NoError(err, "get balance (user_id %d)", id)
	a.EqualValues(change, balance, "compare change and balance")

	err = suite.db.ChangeBalance(ctx, id, -change, "withdraw")
	a.NoError(err, "change balance (user_id %d, change %d)", id, change)

	balance, _, err = suite.db.GetBalance(ctx, id)
	a.NoError(err, "get balance (user_id %d)", id)

	a.EqualValues(0, balance, "balance must be 0")

	err = suite.db.Delete(ctx, id)
	a.NoError(err, "delete user (id %d)", id)
}

func (suite *DBSuite) TestSwitch() {
	ctx := context.Background()

	id, change := rand.Int63n(100000), rand.Int63n(100000)
	a := suite.Assert()

	err := suite.db.ChangeBalance(ctx, id, change, "deposit")
	a.NoError(err, "change balance with creation (user_id %d, change %d)", id, change)

	err = suite.db.Switch(ctx, id, id+1)
	a.NoError(err, "switch (old user_id %d, new %d)", id, id+1)

	balance, _, err := suite.db.GetBalance(ctx, id+1)
	a.NoError(err, "get balance (user_id %d)", id+1)
	a.EqualValues(change, balance, "compare change and balance")

	err = suite.db.Delete(ctx, id+1)
	a.NoError(err, "delete user (id %d)", id+1)
}

func (suite *DBSuite) TestTransfer() {
	ctx := context.Background()

	first, change := rand.Int63n(100000)+1, rand.Int63n(100000)+1
	a := suite.Assert()

	err := suite.db.ChangeBalance(ctx, first, change, "deposit")
	a.NoError(err, "change balance with creation (user_id %d, change %d)", first, change)

	err = suite.db.Transfer(ctx, first, first+1, change, "transfer")
	a.NoError(err, "transfer (from %d,to %d, change %d)", first, first+1, change)

	balance, _, err := suite.db.GetBalance(ctx, first+1)
	a.NoError(err, "get balance (user_id %d)", first+1)
	a.EqualValues(change, balance, "compare change and new balance")

	balance, _, err = suite.db.GetBalance(ctx, first)
	a.NoError(err, "get balance (user_id %d)", first)
	a.EqualValues(0, balance, "from balance must be 0")

	err = suite.db.Delete(ctx, first)
	a.NoError(err, "delete user (id %d)", first)

	err = suite.db.Delete(ctx, first+1)
	a.NoError(err, "delete user (id %d)", first+1)
}

func TestDB(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DBSuite))
}
