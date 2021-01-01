package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateTable creates the users table
func CreateTable(ctx context.Context, pool *pgxpool.Pool) error {
	pgtag, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username VARCHAR(156) NOT NULL UNIQUE,
		hash TEXT NOT NULL
	);`)

	log.Println(pgtag)

	return err
}
