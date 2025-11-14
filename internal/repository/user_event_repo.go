package repository

import (
	"context"
	"errors"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (	ErrNotEnoughSeats = errors.New("event does not have enough seats"))

type UserEventRepository interface {
	CreateUserEvent(ctx context.Context, userEvent *models.UserEvent) error
	GetUserEventInfo(ctx context.Context, id uuid.UUID) (*models.UserEvent, error)
	GetUserEvents(ctx context.Context, userID string, pagination models.Pagination) (*models.PaginatedUserEvents, error)
	UpdateUserEvent(ctx context.Context, id uuid.UUID, newGuestNumber int64, eventId uuid.UUID) (*models.UserEvent, error)
	DeleteUserEvent(ctx context.Context, id uuid.UUID, eventId uuid.UUID) error
}

type userEventRepository struct {
	db *gorm.DB
}

func NewUserEventRepository(db *gorm.DB) UserEventRepository {
	return &userEventRepository{db: db}
}

func (r *userEventRepository) updateEvent(id uuid.UUID, newOccupiedSeats int64) error {
	var event models.Event
	if err := r.db.Where("id_event = ?", id).First(&event).Error; err != nil {
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

func (r *userEventRepository) CreateUserEvent(ctx context.Context, userEvent *models.UserEvent) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.updateEvent(userEvent.EventID, userEvent.NumberOfGuests)
		if err != nil {
			return err
		}
		
		if err := r.db.Create(userEvent).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *userEventRepository) GetUserEventInfo(ctx context.Context, id uuid.UUID) (*models.UserEvent, error) {
	var userEvent *models.UserEvent
	if err := r.db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_event", "event_title", "event_datetime", "little_picture")
	}).Where("id_user_event = ?", id).First(&userEvent).Error; err != nil {
		return nil, err
	}
	return userEvent, nil
}

func (r *userEventRepository) GetUserEvents(ctx context.Context, userID string, pagination models.Pagination) (*models.PaginatedUserEvents, error) {
	var userEvents []models.UserEvent
	var total int64

	if err := r.db.Model(&models.UserEvent{}).Count(&total).Error; err != nil {
    return nil, err
  }

	offset := (pagination.Page - 1) * pagination.PerPage
	
	if err := r.db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_event", "event_title", "event_datetime", "little_picture")
	}).Where("user_id = ?", userID).Limit(pagination.PerPage).Offset(offset).Find(&userEvents).Error; err != nil {
		return nil, err
	}

	return &models.PaginatedUserEvents{
		UserEvents: userEvents,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	}, nil
}

func (r *userEventRepository) UpdateUserEvent(ctx context.Context, id uuid.UUID, newGuestNumber int64, eventId uuid.UUID) (*models.UserEvent, error) {
	userEvent, err := r.GetUserEventInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = r.updateEvent(userEvent.EventID, newGuestNumber-userEvent.NumberOfGuests)
		if err != nil {
			return err
		}
		userEvent.NumberOfGuests = newGuestNumber
		if err := r.db.Save(userEvent).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return userEvent, nil
}

func (r *userEventRepository) DeleteUserEvent(ctx context.Context, id uuid.UUID, eventId uuid.UUID) error {
	userEvent, err := r.GetUserEventInfo(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = r.updateEvent(eventId, -userEvent.NumberOfGuests)
		if err != nil {
			return err
		}

		if err := r.db.Where("id_user_event = ?", id).Delete(&models.UserEvent{}).Error; err != nil {
			return err
		}
		return nil
	})
	
	if err != nil {
		return err
	}
	return nil
}