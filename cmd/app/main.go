package main

import (
	"context"
	"log"
	"time"

	"github.com/TaperoOO5536/special_backend/internal/app"
	"github.com/TaperoOO5536/special_backend/internal/config"
)

func main() {
	config.LoadEnv()

	cfg := &app.Config{
		Port:         "8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Dsn:          config.GetDsn(),
	}

	app := app.New(cfg)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}