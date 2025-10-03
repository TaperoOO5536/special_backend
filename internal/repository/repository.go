package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Item      ItemRepository
	Event     EventRepository
	User      UserRepository
	Order     OrderRepository
	UserEvent UserEventRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:      NewItemRepository(db),
		Event:     NewEventRepository(db),
		User:      NewUserRepository(db),
		Order:     NewOrderRepository(db),
		UserEvent: NewUserEventRepository(db),
	}
}