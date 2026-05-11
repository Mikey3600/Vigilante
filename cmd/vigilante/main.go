package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "vigilante", Short: "Vigilante observability platform"}

func init() { _ = godotenv.Load() }

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("command_failed: %v", err)
		os.Exit(1)
	}
}
