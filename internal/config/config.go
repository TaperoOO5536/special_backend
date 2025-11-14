package config

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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