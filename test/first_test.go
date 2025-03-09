package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kvloginov/pg-isolation-levels-demo/test/tools"
)

func TestOne(t *testing.T) {
	db := tools.ConnectToDB(t)
	ctx := context.Background()
	tools.MigrateTestDB(t, db)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	require.NoError(t, err, "begin transaction")

	_, err = tx.Exec("INSERT INTO wallets (user_id, balance) VALUES (1, 100)")
	require.NoError(t, err, "insert a row into wallets")

	row := tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = 1")
	var balance int
	err = row.Scan(&balance)
	require.NoError(t, err, "scan balance")
	require.Equal(t, 100, balance, "balance should be 100")

	err = tx.Commit()
	require.NoError(t, err, "commit transaction")
}
