package db

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
