package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "vigilante", Short: "Vigilante observability platform"}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP and gRPC servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunServe(context.Background())
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database schema migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunMigrate(context.Background())
	},
}

func init() {
	_ = godotenv.Load()
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
