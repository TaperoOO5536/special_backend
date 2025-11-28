package config

import (
	"database/sql"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	Client yookassa.Client
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

func NewYookassaClient(shopId, secretKey string) *Client {
	client := yookassa.NewClient(shopId, secretKey)
	return &Client{
		Client: *client,
	}
}