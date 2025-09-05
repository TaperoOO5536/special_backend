package config

import (
	// "context"
	// "time"

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

func GetEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewDBClient(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic("failed to get database instance")
	// }
	// defer sqlDB.Close()
	return db
}