package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type ItemRepository interface {
	GetItemInfo(ctx context.Context, id uuid.UUID) (*models.Item, error)
	GetItems(ctx context.Context, pagination models.Pagination) (*models.PaginatedItems, error)
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

func (r *itemRepository) GetItems(ctx context.Context, pagination models.Pagination) (*models.PaginatedItems, error) {
	var items []models.Item
	var total int64

	if err := r.db.Model(&models.Item{}).Count(&total).Error; err != nil {
    return nil, err
  }

	offset := (pagination.Page - 1) * pagination.PerPage
	if err := r.db.Limit(pagination.PerPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	return &models.PaginatedItems{
		Items:      items,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	}, nil
}