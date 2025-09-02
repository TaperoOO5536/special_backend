package models

import (
	"github.com/google/uuid"
	// "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	ID          string      `gorm:"column:ID_User"`
	Name        string      `gorm:"column:F_N_User"`
	Surname     string      `gorm:"column:S_N_User"`
	Nickname    string      `gorm:"column:N_N_User"`
	PhoneNumber string      `gorm:"column:Phone_N_User"`
	UserIvents  []UserIvent `gorm:"foreignKey:UserID"`
	Orders      []Order     `gorm:"foreignKey:UserID"`
}

type Ivent struct {
	ID            uuid.UUID      `gorm:"column:ID_Ivent"`
	Title         string         `gorm:"column:Ivent_Title"`
	Description   string         `gorm:"column:Ivent_Description"`
	DateTime      time.Time      `gorm:"column:Ivent_DateTime"`
	Price         int64          `gorm:"column:Ivent_Price"`
	TotalSeats    int64          `gorm:"column:Total_Seats"`
	OccupiedSeats int64          `gorm:"column:Occupied_Seats"`
	LittlePicture string         `gorm:"column:Little_Picture"`
	UserIvents    []UserIvent    `gorm:"foreignKey:IventID"`
	Pictures      []IventPicture `gorm:"foreignKey:IventID"`
}

type UserIvent struct {
	ID             uuid.UUID `gorm:"column:ID_User_Ivent"`
	UserID         string    `gorm:"column:User_ID;index"`
	IventID        uuid.UUID `gorm:"column:Ivent_ID;index"`
	NumberOfGuests int64     `gorm:"column:Number_Of_Guests"`
}

type IventPicture struct {
	ID      uuid.UUID `gorm:"column:ID_Ivent_Picture"`
	IventID uuid.UUID `gorm:"column:Ivent_ID"`
	Path    string    `gorm:"column:Picture_Path"`
}

type Order struct {
	ID             uuid.UUID   `gorm:"column:ID_Order"`
	Number         string      `gorm:"column:Order_Number"`
	UserID         string      `gorm:"column:User_ID"`
	FormDate       time.Time   `gorm:"column:Order_Form_DateTime"`
	CompletionDate time.Time   `gorm:"column:Completion_Date"`
	Comment        string      `gorm:"column:Order_Comment"`
	Status         string      `gorm:"column:Order_Status"`
	OrderAmount    int64       `gorm:"column:Order_Amount"`
	OrderItems     []OrderItem `gorm:"foreignKey:OrderID"`
}

type Item struct {
	ID            uuid.UUID     `gorm:"column:ID_Item"`
	Title         string        `gorm:"column:Item_Title"`
	Description   string        `gorm:"column:Item_Description"`
	Price         int64         `gorm:"column:Item_Price"`
	LittlePicture string        `gorm:"column:Little_Picture"`
	OrderItems    []OrderItem   `gorm:"foreignKey:ItemID"`
	Pictures      []ItemPicture `gorm:"foreignKey:ItemID"`
}

type OrderItem struct{
	ID       uuid.UUID `gorm:"column:ID_Order_Item"`
	OrderID  uuid.UUID `gorm:"column:Order_ID;index"`
	ItemID   uuid.UUID `gorm:"column:Item_ID;index"`
	Quantity int64
}

type ItemPicture struct {
	ID     uuid.UUID `gorm:"column:ID_Item_Picture"`
	ItemID uuid.UUID `gorm:"column:Item_ID"`
	Path   string    `gorm:"column:Picture_Path"`
}