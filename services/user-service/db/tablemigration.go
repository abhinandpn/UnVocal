package db

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

var CreateUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    number VARCHAR(20) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func GetTableQueries() []string {
	return []string{
		CreateUsersTable, // users table
	}
}

func RunMigrations(db *pgxpool.Pool) error {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		query, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(context.Background(), string(query)); err != nil {
			return err
		}
	}

	return nil
}
