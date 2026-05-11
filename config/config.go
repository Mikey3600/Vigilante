package config

import (
	"fmt"
	"os"
)

type Config struct { DatabaseURL, JWTSecret, GRPCPort, HTTPPort, GeminiAPIKey, SlackWebhookURL, SMTPHost, SMTPPort, SMTPUser, SMTPPass, AlertFromEmail, Env, AllowedOrigins string }

func Load() (Config, error) {
	cfg := Config{DatabaseURL: os.Getenv("DATABASE_URL"), JWTSecret: os.Getenv("JWT_SECRET"), GRPCPort: getDefault("GRPC_PORT", "50051"), HTTPPort: getDefault("HTTP_PORT", "8080"), GeminiAPIKey: os.Getenv("GEMINI_API_KEY"), SlackWebhookURL: os.Getenv("SLACK_WEBHOOK_URL"), SMTPHost: os.Getenv("SMTP_HOST"), SMTPPort: os.Getenv("SMTP_PORT"), SMTPUser: os.Getenv("SMTP_USER"), SMTPPass: os.Getenv("SMTP_PASS"), AlertFromEmail: os.Getenv("ALERT_FROM_EMAIL"), Env: getDefault("ENV", "development"), AllowedOrigins: getDefault("ALLOWED_ORIGINS", "*")}
	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" { return Config{}, fmt.Errorf("missing required env vars: DATABASE_URL and JWT_SECRET") }
	return cfg, nil
}
func getDefault(key, fallback string) string { if v := os.Getenv(key); v != "" { return v }; return fallback }
