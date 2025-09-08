// заглушка пока

package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (error)
	GetUserInfo(ctx context.Context) (*models.User, error)
	UpdateUser(ctx context.Context, updates *UserUpdates) (*models.User, error)
}

type UserUpdates struct {
	Name        string
	Surname     string
	Nickname    string
	PhoneNumber string
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (error) {
	
	
	return nil
}

func (r *userRepository) GetUserInfo(ctx context.Context) (*models.User, error) {
	return nil, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, updates *UserUpdates) (*models.User, error) {
	return nil, nil
}
