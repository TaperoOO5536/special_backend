package service

import (
	"context"
	"errors"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound =    errors.New("user not found")
	ErrInvalidInitData = errors.New("invalid init data")
	ErrUserExists =      errors.New("user already exists")
)

type UserService struct {
	userRepo repository.UserRepository
	token    string
}

func NewUserService(userRepo repository.UserRepository, token string) *UserService {
	return &UserService{
		userRepo: userRepo,
		token: token,
	}
}

func (s *UserService) CreateUser(ctx context.Context, initData string) (error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	existingUser, _ := s.userRepo.GetUserInfo(ctx, user.ID)
	if existingUser != nil {
		return ErrUserExists
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *UserService) GetUserInfo(ctx context.Context, initData string) (*models.User, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.userRepo.GetUserInfo(ctx, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return userInfo, nil
}

func (s *UserService) UpdateUser(ctx context.Context, initData string, phoneNumber string) (*models.User, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetUserInfo(ctx, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if phoneNumber != "" {
		user.PhoneNumber = phoneNumber
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return updatedUser, nil
}