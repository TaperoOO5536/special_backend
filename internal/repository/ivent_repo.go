package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type IventRepository interface {
	GetIventInfo(ctx context.Context, id uuid.UUID) (*models.Ivent, error)
	GetIvents(ctx context.Context) ([]*models.Ivent, error)
}

type iventRepository struct {
	db *gorm.DB
}

func NewIventRepository(db *gorm.DB) IventRepository {
	return &iventRepository{db: db}
}

func (r *iventRepository) GetIventInfo(ctx context.Context, id uuid.UUID) (*models.Ivent, error) {
	var ivent models.Ivent
	if err := r.db.Preload("Pictures").Where("ID_Ivent = ?", id).First(&ivent).Error; err != nil {
		return nil, err
	}
	return &ivent, nil
}

func (r *iventRepository) GetIvents(ctx context.Context) ([]*models.Ivent, error) {
	var ivents []*models.Ivent
	if err := r.db.Find(&ivents).Error; err != nil {
		return nil, err
	}
	return ivents, nil
}