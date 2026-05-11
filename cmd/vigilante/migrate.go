package main

import (
	"context"
	"fmt"
	"os"

	"github.com/user/vigilante/internal/storage"
)

func RunMigrate(ctx context.Context) error {
	fmt.Println("Vigilante starting...")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	db, err := storage.NewDB(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	fmt.Println("Connected to database")

	if err := db.RunMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("Migrations complete")

	return nil
}
