package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(databaseURL string) error {
	var err error

	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return err
	}

	if err := DB.Ping(context.Background()); err != nil {
		return err
	}

	// Run migrations
	if err := TableMigrate(DB); err != nil {
		return err
	}

	log.Println("✅ PostgreSQL connected")
	return nil
}
func TableMigrate(pool *pgxpool.Pool) error {
	queries := GetTableQueries()

	for _, query := range queries {
		if _, err := pool.Exec(context.Background(), query); err != nil {
			return err
		}
	}

	log.Println("✅ Database table migration completed successfully.")

	return nil
}
