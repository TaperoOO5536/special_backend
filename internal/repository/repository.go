package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Item  ItemRepository
	Ivent IventRepository
	User  UserRepository
	Order OrderRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:  NewItemRepository(db),
		Ivent: NewIventRepository(db),
		User:  NewUserRepository(db),
		Order: NewOrderRepository(db),
	}
}