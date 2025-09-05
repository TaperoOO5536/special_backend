package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type ItemRepository interface {
	GetItemInfo(ctx context.Context, id uuid.UUID) (*models.Item, error)
	GetItems(ctx context.Context) ([]*models.Item, error)
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) GetItemInfo(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	var item models.Item
	if err := r.db.Preload("Pictures").Where("id_item = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) GetItems(ctx context.Context) ([]*models.Item, error) {
	var items []*models.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}