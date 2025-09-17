package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserIventRepository interface {
	CreateUserIvent(ctx context.Context, userIvent *models.UserIvent) error
	GetUserIventInfo(ctx context.Context, id uuid.UUID) (*models.UserIvent, error)
	GetUserIvents(ctx context.Context, userID string) ([]*models.UserIvent, error)
	UpdateUserIvent(ctx context.Context, id uuid.UUID, newGuestNumber int64) (*models.UserIvent, error)
	DeleteUserIvent(ctx context.Context, id uuid.UUID) error
}

type userIventRepository struct {
	db *gorm.DB
}

func NewUserIventRepository(db *gorm.DB) UserIventRepository {
	return &userIventRepository{db: db}
}

func (r *userIventRepository) CreateUserIvent(ctx context.Context, userIvent *models.UserIvent) error {
	if err := r.db.Create(userIvent).Error; err != nil {
		return err
	}
	return nil
}

func (r *userIventRepository) GetUserIventInfo(ctx context.Context, id uuid.UUID) (*models.UserIvent, error) {
	var userIvent *models.UserIvent
	if err := r.db.Preload("Ivent", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_ivent", "ivent_title", "ivent_datetime", "little_picture", "mime_type")
	}).Where("id_user_ivent = ?", id).First(&userIvent).Error; err != nil {
		return nil, err
	}
	return userIvent, nil
}

func (r *userIventRepository) GetUserIvents(ctx context.Context, userID string) ([]*models.UserIvent, error) {
	var userIvents []*models.UserIvent
	if err := r.db.Preload("Ivent", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_ivent", "ivent_title", "ivent_datetime", "little_picture", "mime_type")
	}).Where("user_id = ?", userID).Find(&userIvents).Error; err != nil {
		return nil, err
	}

	return userIvents, nil
}

func (r *userIventRepository) UpdateUserIvent(ctx context.Context, id uuid.UUID, newGuestNumber int64) (*models.UserIvent, error) {
	userIvent, err := r.GetUserIventInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	userIvent.NumberOfGuests = newGuestNumber
	if err := r.db.Save(userIvent).Error; err != nil {
		return nil, err
	}
	return userIvent, nil
}

func (r *userIventRepository) DeleteUserIvent(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Where("id_user_ivent = ?", id).Delete(&models.UserIvent{}).Error; err != nil {
		return err
	}

	return nil
}