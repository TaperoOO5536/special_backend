package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Item      ItemRepository
	Ivent     IventRepository
	User      UserRepository
	Order     OrderRepository
	UserIvent UserIventRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:      NewItemRepository(db),
		Ivent:     NewIventRepository(db),
		User:      NewUserRepository(db),
		Order:     NewOrderRepository(db),
		UserIvent: NewUserIventRepository(db),
	}
}