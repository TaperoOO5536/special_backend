package repository

import (
	"context"
	"errors"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

var (	ErrNotEnoughSeats = errors.New("event does not have enough seats"))

type EventRepository interface {
	GetEventInfo(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEvents(ctx context.Context) ([]*models.Event, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, newOccupiedSeats int64) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetEventInfo(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	var event models.Event
	if err := r.db.Preload("Pictures").Where("id_event = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetEvents(ctx context.Context) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) UpdateEvent(ctx context.Context, id uuid.UUID, newOccupiedSeats int64) error {
	event, err := r.GetEventInfo(ctx, id)
	if err != nil {
		return err
	}

	if (event.TotalSeats - (event.OccupiedSeats + newOccupiedSeats) < 0) {
		return ErrNotEnoughSeats
	}

	event.OccupiedSeats = event.OccupiedSeats + newOccupiedSeats
	if err := r.db.Save(event).Error; err != nil {
		return err
	}

	return nil
}