package tools

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/kvloginov/pg-isolation-levels-demo/internal/infra/pg"
)

func MigrateTestDB(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			user_id INT PRIMARY KEY,
			balance INT
		);
	`)
	require.NoError(t, err, "create test table")

}

func ConnectToDB(t *testing.T) *sqlx.DB {
	db, err := pg.NewDB(defaultDBConfig())
	require.NoError(t, err, "connect to database")
	return db
}

func defaultDBConfig() pg.Config {
	return pg.Config{
		Database:    "demo",
		Username:    "user",
		Password:    "pass",
		HostPrimary: "localhost",
		Port:        "25432",
	}
}
