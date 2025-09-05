package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Item  ItemRepository
	Ivent IventRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:  NewItemRepository(db),
		Ivent: NewIventRepository(db),
	}
}