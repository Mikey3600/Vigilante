package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vigilante/internal/storage"
)

var migrateCmd = &cobra.Command{Use: "migrate", Short: "Runs database schema migrations", RunE: func(cmd *cobra.Command, args []string) error {
	db, err := storage.NewDB(context.Background(), os.Getenv("DATABASE_URL")); if err != nil { return err }
	defer db.Close(); if err := db.RunMigrations(context.Background()); err != nil { return err }
	slog.Info("migrations_applied")
	return nil
}}

func init() { rootCmd.AddCommand(migrateCmd) }
