package config

import (
	"os"

	"github.com/elton/cerp-sync/utils/logger"
	"github.com/joho/godotenv"
)

// Config returns a Config object by specified key.
func Config(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		logger.Error.Printf("Error loading .env file: %s", err)
	}
	return os.Getenv(key)
}
