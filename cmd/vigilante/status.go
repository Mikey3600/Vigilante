package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the runtime status of the API",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:3000/health")
		if err != nil {
			fmt.Printf("Server unreachable: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			fmt.Println("Status: OK (200)")
		} else {
			fmt.Printf("Status Failed: (%d)\n", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
