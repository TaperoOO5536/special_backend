package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	databaseURL := os.Getenv("DSN")
	if databaseURL == "" {
		return fmt.Errorf("DSN is not set")
	}

	return nil
}

func GetDatabaseURL() string {
	return os.Getenv("DSN")
}



