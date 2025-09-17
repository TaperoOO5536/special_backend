package repository

import (
	"context"
	"errors"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

var (	ErrNotEnoughSeats = errors.New("ivent does not have enough seats"))

type IventRepository interface {
	GetIventInfo(ctx context.Context, id uuid.UUID) (*models.Ivent, error)
	GetIvents(ctx context.Context) ([]*models.Ivent, error)
	UpdateIvent(ctx context.Context, id uuid.UUID, newOccupiedSeats int64) error
}

type iventRepository struct {
	db *gorm.DB
}

func NewIventRepository(db *gorm.DB) IventRepository {
	return &iventRepository{db: db}
}

func (r *iventRepository) GetIventInfo(ctx context.Context, id uuid.UUID) (*models.Ivent, error) {
	var ivent models.Ivent
	if err := r.db.Preload("Pictures").Where("id_ivent = ?", id).First(&ivent).Error; err != nil {
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

func (r *iventRepository) UpdateIvent(ctx context.Context, id uuid.UUID, newOccupiedSeats int64) error {
	ivent, err := r.GetIventInfo(ctx, id)
	if err != nil {
		return err
	}

	if (ivent.TotalSeats - (ivent.OccupiedSeats + newOccupiedSeats) < 0) {
		return ErrNotEnoughSeats
	}

	ivent.OccupiedSeats = ivent.OccupiedSeats + newOccupiedSeats
	if err := r.db.Save(ivent).Error; err != nil {
		return err
	}

	return nil
}