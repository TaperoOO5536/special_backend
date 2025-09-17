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
	ErrIventNotFound = errors.New("ivent not found")
)

type IventService struct {
	iventRepo repository.IventRepository
}

func NewIventService(iventRepo repository.IventRepository) *IventService {
	return &IventService{
		iventRepo: iventRepo,
	}
}

func (s *IventService) GetIventInfo(ctx context.Context, id uuid.UUID) (*models.Ivent, error) {
	ivent, err := s.iventRepo.GetIventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrIventNotFound
		}
		return nil, err
	}
	
	return ivent, nil
}

func (s *IventService) GetIvents(ctx context.Context) ([]*models.Ivent, error) {
	ivents, err := s.iventRepo.GetIvents(ctx)
	if err != nil {
		return nil, err
	}

	return ivents, nil
}