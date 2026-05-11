package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/user/vigilante/internal/auth"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate a JWT token",
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return fmt.Errorf("JWT_SECRET is not set")
		}

		token, err := auth.GenerateToken(secret, "default", "", 24*time.Hour)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		fmt.Println(token)
		return nil
	},
}
