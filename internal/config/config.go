package config

import (
	// "context"
	// "time"

	// "github.com/jackc/pgx/v5"
	// "github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDBClient(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}
	defer sqlDB.Close()
	return db
}