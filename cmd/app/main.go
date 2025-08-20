package main

import (
	// "context"
	"log"
	// "time"

	// "github.com/TaperoOO5536/special_backend/internal/config"
	"github.com/TaperoOO5536/special_backend/pkg/env"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	// db := config.newDBClient(env.GetDatabaseURL())
}