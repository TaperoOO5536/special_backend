package service

import (
	"context"
	"errors"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type ItemService struct {
	itemRepo repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

func (s *ItemService) GetItemInfo(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item, err := s.itemRepo.GetItemInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	
	return item, nil
}

func (s *ItemService) GetItems(ctx context.Context, pagination models.Pagination) (*models.PaginatedItems, error) {
	items, err := s.itemRepo.GetItems(ctx, pagination)
	if err != nil {
		return nil, err
	}
	
	return items, nil
}