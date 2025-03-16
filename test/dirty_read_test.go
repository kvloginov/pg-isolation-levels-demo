package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kvloginov/pg-isolation-levels-demo/test/tools"
)

func TestDirtyRead(t *testing.T) {
	if !tools.IsPostgresql() { // Postgresql doesn't support ReadUncommitted isolation level
		balance, err := TryReadFromNotCommitedTransaction(t, sql.LevelReadUncommitted, sql.LevelReadUncommitted)
		require.NoError(t, err)
		// with ReadUncommitted isolation level we can read uncommited data
		require.Equal(t, 100, balance)
	}
	_, err := TryReadFromNotCommitedTransaction(t, sql.LevelReadCommitted, sql.LevelReadCommitted)
	// with ReadCommitted isolation level we can't read uncommited data
	require.ErrorContains(t, err, sql.ErrNoRows.Error())
}

func TryReadFromNotCommitedTransaction(t *testing.T, tx1Isolation, tx2Isolation sql.IsolationLevel) (int, error) {
	db := tools.ConnectToDB(t)
	ctx := context.Background()
	tools.DropTestDB(t, db)
	tools.MigrateTestDB(t, db)

	// Begin transaction 1
	tx1, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: tx1Isolation})
	require.NoError(t, err, "begin transaction 1")

	// Begin transaction 2
	tx2, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: tx2Isolation})
	require.NoError(t, err, "begin transaction 2")

	// insert a row into wallet
	_, err = tx1.Exec("INSERT INTO wallets (user_id, balance) VALUES (1, 100)")
	require.NoError(t, err, "insert a row into wallets")

	// select balance from NOT COMMITED wallet in transaction 1
	row := tx2.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = 1")
	var balance int
	err = row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("scan balance: %w", err)
	}

	err = tx1.Commit()
	require.NoError(t, err, "commit transaction 1")

	err = tx2.Rollback()
	require.NoError(t, err, "rollback transaction 2")

	// return balance, that was read from not commited transaction
	return balance, nil
}
