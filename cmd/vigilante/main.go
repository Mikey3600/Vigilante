package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "vigilante",
	Short: "Vigilante is a Backend Observability & Incident Intelligence Platform",
	Long:  `Vigilante ingests logs and metrics securely, stores them in TimescaleDB, detects anomalies, and uses Gemini AI to give root calls via Webhook logic bounds.`,
}

func init() {
	_ = godotenv.Load()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
