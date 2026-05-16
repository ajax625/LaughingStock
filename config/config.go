package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	LaughingStockBaseURL string

	TradeXBaseURL string
	TradeXAPIKey  string

	TelegramBotToken string

	AllowedOrigins string
}

func Load() *Config {
	return &Config{
		Port:                 getEnvOrDefault("PORT", "9090"),
		DatabaseURL:          getEnvOrDefault("DATABASE_URL", "postgres://laughingstock:laughingstock@localhost:5433/laughingstock?sslmode=disable"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		LaughingStockBaseURL: getEnvOrDefault("LAUGHINGSTOCK_BASE_URL", "https://laughingstock.informatrixlabs.com"),
		TradeXBaseURL:        getEnvOrDefault("TRADEX_BASE_URL", "https://tradex.informatrixlabs.com"),
		TradeXAPIKey:         os.Getenv("TRADEX_API_KEY"),
		TelegramBotToken:     os.Getenv("TELEGRAM_BOT_TOKEN"),
		AllowedOrigins:       getEnvOrDefault("ALLOWED_ORIGINS", "*"),
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var _ = strconv.Itoa // suppress unused import
