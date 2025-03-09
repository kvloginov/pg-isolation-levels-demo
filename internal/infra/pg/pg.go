package pg

import (
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Database    string
	Username    string
	Password    string
	HostPrimary string
	Port        string
}

func NewDB(cfg Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.HostPrimary, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	db, err := sqlx.Connect("pgx", dataSource)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return db, nil
}
