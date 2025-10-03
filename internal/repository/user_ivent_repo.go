package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEventRepository interface {
	CreateUserEvent(ctx context.Context, userEvent *models.UserEvent) error
	GetUserEventInfo(ctx context.Context, id uuid.UUID) (*models.UserEvent, error)
	GetUserEvents(ctx context.Context, userID string) ([]*models.UserEvent, error)
	UpdateUserEvent(ctx context.Context, id uuid.UUID, newGuestNumber int64) (*models.UserEvent, error)
	DeleteUserEvent(ctx context.Context, id uuid.UUID) error
}

type userEventRepository struct {
	db *gorm.DB
}

func NewUserEventRepository(db *gorm.DB) UserEventRepository {
	return &userEventRepository{db: db}
}

func (r *userEventRepository) CreateUserEvent(ctx context.Context, userEvent *models.UserEvent) error {
	if err := r.db.Create(userEvent).Error; err != nil {
		return err
	}
	return nil
}

func (r *userEventRepository) GetUserEventInfo(ctx context.Context, id uuid.UUID) (*models.UserEvent, error) {
	var userEvent *models.UserEvent
	if err := r.db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_event", "event_title", "event_datetime", "little_picture", "mime_type")
	}).Where("id_user_event = ?", id).First(&userEvent).Error; err != nil {
		return nil, err
	}
	return userEvent, nil
}

func (r *userEventRepository) GetUserEvents(ctx context.Context, userID string) ([]*models.UserEvent, error) {
	var userEvents []*models.UserEvent
	if err := r.db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_event", "event_title", "event_datetime", "little_picture", "mime_type")
	}).Where("user_id = ?", userID).Find(&userEvents).Error; err != nil {
		return nil, err
	}

	return userEvents, nil
}

func (r *userEventRepository) UpdateUserEvent(ctx context.Context, id uuid.UUID, newGuestNumber int64) (*models.UserEvent, error) {
	userEvent, err := r.GetUserEventInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	userEvent.NumberOfGuests = newGuestNumber
	if err := r.db.Save(userEvent).Error; err != nil {
		return nil, err
	}
	return userEvent, nil
}

func (r *userEventRepository) DeleteUserEvent(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Where("id_user_event = ?", id).Delete(&models.UserEvent{}).Error; err != nil {
		return err
	}

	return nil
}