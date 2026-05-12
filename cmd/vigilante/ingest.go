package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var inputFile string

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Sends a log payload file to the running HTTP API",
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			log.Fatal("Must provide --file flag")
		}

		data, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Failed reading file: %v", err)
		}

		// Simplified representation for testing locally
		resp, err := http.Post("http://localhost:3000/api/v1/logs", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Fatalf("Failed to request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, string(body))
	},
}

func init() {
	ingestCmd.Flags().StringVar(&inputFile, "file", "", "Path to the JSON file to ingest")
	rootCmd.AddCommand(ingestCmd)
}
