package models

import (
	"time"
)

type CreatedOrder struct {
	Number         int32		 `json:"number"`
	UserID         string    `json:"user_id"`
	CompletionDate time.Time `json:"completion_date"`
	OrderAmount    int64     `json:"order_amount"`
}