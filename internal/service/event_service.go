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
	ErrEventNotFound = errors.New("event not found")
)

type EventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (s *EventService) GetEventInfo(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event, err := s.eventRepo.GetEventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	
	return event, nil
}

func (s *EventService) GetEvents(ctx context.Context, pagination models.Pagination) (*models.PaginatedEvents, error) {
	events, err := s.eventRepo.GetEvents(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return events, nil
}