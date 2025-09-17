package models

import (
	"github.com/google/uuid"
	// "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	ID          string      `gorm:"column:id_user;primaryKey"`
	Name        string      `gorm:"column:f_n_user"`
	Surname     string      `gorm:"column:s_n_user"`
	Nickname    string      `gorm:"column:n_n_user"`
	PhoneNumber string      `gorm:"column:phone_n_user"`
	UserIvents  []UserIvent `gorm:"foreignKey:UserID"`
	Orders      []Order     `gorm:"foreignKey:UserID"`
}

type Ivent struct {
	ID            uuid.UUID      `gorm:"column:id_ivent;primaryKey"`
	Title         string         `gorm:"column:ivent_title"`
	Description   string         `gorm:"column:ivent_description"`
	DateTime      time.Time      `gorm:"column:ivent_datetime"`
	Price         int64          `gorm:"column:ivent_price"`
	TotalSeats    int64          `gorm:"column:total_seats"`
	OccupiedSeats int64          `gorm:"column:occupied_seats"`
	LittlePicture []byte         `gorm:"column:little_picture"`
	MimeType			string         `gorm:"column:mime_type"`
	UserIvents    []UserIvent    `gorm:"foreignKey:IventID"`
	Pictures      []IventPicture `gorm:"foreignKey:IventID"`
}

type UserIvent struct {
	ID             uuid.UUID `gorm:"column:id_user_ivent;primaryKey"`
	UserID         string    `gorm:"column:user_id"`
	IventID        uuid.UUID `gorm:"column:ivent_id"`
	NumberOfGuests int64     `gorm:"column:number_of_guests"`
	Ivent          Ivent     `gorm:"foreignKey:IventID"`
}

type IventPicture struct {
	ID      uuid.UUID `gorm:"column:id_ivent_picture;primaryKey"`
	IventID uuid.UUID `gorm:"column:ivent_id"`
	Path    []byte    `gorm:"column:picture_path"`
	MimeType string   `gorm:"column:mime_type"`
}

type Order struct {
	ID             uuid.UUID   `gorm:"column:id_order;primaryKey"`
	Number         string      `gorm:"column:order_number"`
	UserID         string      `gorm:"column:user_id"`
	FormDate       time.Time   `gorm:"column:order_form_datetime"`
	CompletionDate time.Time   `gorm:"column:completion_date"`
	Comment        string      `gorm:"column:order_comment"`
	Status         string      `gorm:"column:order_status"`
	OrderAmount    int64       `gorm:"column:order_amount"`
	OrderItems     []OrderItem `gorm:"foreignKey:OrderID"`
}

type Item struct {
	ID            uuid.UUID     `gorm:"column:id_item;primaryKey"`
	Title         string        `gorm:"column:item_title"`
	Description   string        `gorm:"column:item_description"`
	Price         int64         `gorm:"column:item_price"`
	LittlePicture []byte        `gorm:"column:little_picture"`
	MimeType			string        `gorm:"column:mime_type"`
	Pictures      []ItemPicture `gorm:"foreignKey:ItemID"`
}

type OrderItem struct{
	ID       uuid.UUID `gorm:"column:id_order_item;primaryKey"`
	OrderID  uuid.UUID `gorm:"column:order_id;index"`
	ItemID   uuid.UUID `gorm:"column:item_id;index"`
	Quantity int64
	Item 	   Item      `gorm:"foreignKey:ItemID"`
}

type ItemPicture struct {
	ID     uuid.UUID `gorm:"column:id_item_picture;primaryKey"`
	ItemID uuid.UUID `gorm:"column:item_id"`
	Path   []byte    `gorm:"column:picture_path"`
	MimeType string  `gorm:"column:mime_type"`
}