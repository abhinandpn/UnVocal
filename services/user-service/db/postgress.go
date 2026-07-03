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

	log.Println("✅ PostgreSQL connected")

	// Run migrations
	if err := RunMigrations(DB); err != nil {
		return err
	}

	log.Println("✅ Database migrations completed successfully.")

	return nil
}
