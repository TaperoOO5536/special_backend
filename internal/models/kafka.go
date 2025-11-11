package models

import (
	"time"
)

type KafkaOrder struct {
	Number         string		 `json:"number"`
	UserID         string    `json:"user_id"`
	CompletionDate time.Time `json:"completion_date"`
	OrderAmount    int64     `json:"order_amount"`
}

type KafkaUserEvent struct {
	UserNickName       string `json:"user_nickname"`
	EventTitle         string `json:"event_title"`
	EventOccupiedSeats int64  `json:"event_occupied_seats"`
	EventTotalSeats    int64  `json:"event_total_seats"`
	NumberOfGuests     int64  `json:"number_of_guests"`
}