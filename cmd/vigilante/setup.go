package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/spf13/cobra"
    "github.com/user/vigilante/internal/storage"
)

var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Initialize Vigilante with a default tenant, service, and JWT token",
    RunE: func(cmd *cobra.Command, args []string) error {
        ctx := context.Background()
        db, err := storage.NewDB(ctx, os.Getenv("DATABASE_URL"))
        if err != nil {
            return fmt.Errorf("failed to connect to database: %w", err)
        }
        defer db.Close()
        tenantID := "22222222-2222-2222-2222-222222222222"
        serviceID := "11111111-1111-1111-1111-111111111111"
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
            jwtSecret = "vigilante-local-secret-key"
        }
        db.Pool.Exec(ctx, `INSERT INTO tenants (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, tenantID, "default")
        db.Pool.Exec(ctx, `INSERT INTO services (id, tenant_id, name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, serviceID, tenantID, "api-gateway")
        claims := jwt.MapClaims{
            "tenant_id": tenantID,
            "user_id":   "admin",
            "exp":       time.Now().Add(24 * time.Hour).Unix(),
            "iat":       time.Now().Unix(),
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        signed, _ := token.SignedString([]byte(jwtSecret))
        fmt.Println("Setup complete!")
        fmt.Println("Tenant ID:  " + tenantID)
        fmt.Println("Service ID: " + serviceID)
        fmt.Println("JWT Token:  " + signed)
        fmt.Println("Open http://localhost:3000 and paste the token.")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(setupCmd)
}
