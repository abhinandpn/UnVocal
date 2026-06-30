package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(databaseURL string) error {

	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("Error connecting to PostgreSQL:", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		fmt.Println("Error pinging PostgreSQL:", err)
	}
	log.Println("\n✅ PostgreSQL connected")

	return nil
}
func TableMigrate(pool *pgxpool.Pool) error {
	queries := GetTableQueries()

	for _, query := range queries {
		if _, err := pool.Exec(context.Background(), query); err != nil {
			return err
		}
	}

	fmt.Println("✅ Table migration completed successfully.")
	return nil
}
