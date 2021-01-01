package db

import (
	"context"
	"fmt"

	"davidws/config"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PGXConnect connects to a Postgres instance and returns a *pgx.ConnPool and an error
func PGXConnect(ctx context.Context) (*pgxpool.Pool, error) {
	if config.DBURI != "" {
		return pgxpool.Connect(ctx, config.DBURI)
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=10",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	//conf, err := pgx.ParseConfig(connectionString)
	_, err := pgx.ParseConfig(connectionString)

	if err != nil {
		return nil, err
	}

	return pgxpool.Connect(ctx, connectionString)
}
