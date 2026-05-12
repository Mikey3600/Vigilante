package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vigilante/internal/storage"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database schema migrations",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := os.Getenv("DATABASE_URL")
		db, err := storage.NewDB(context.Background(), dbUrl)
		if err != nil {
			log.Fatalf("Failed to init db: %v", err)
		}
		defer db.Close()

		if err := db.RunMigrations(context.Background()); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations applied successfully!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
