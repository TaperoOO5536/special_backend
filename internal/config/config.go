package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	err := godotenv.Load("pkg/env/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetDsn() string {
	return GetEnvDefault("DSN", "")
}

func GetToken() string {
	return GetEnvDefault("TOKEN", "")
}

func GetEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewDBClient(dsn string) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}
	return db, sqlDB
}