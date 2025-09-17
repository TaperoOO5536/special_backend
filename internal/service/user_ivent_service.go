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
	ErrUserIventNotFound = errors.New("userivent not found")
)

type UserIventService struct {
	iventRepo     repository.IventRepository
	userIventRepo repository.UserIventRepository
	token         string
}

func NewUserIventService(userIventRepo repository.UserIventRepository, iventRepo repository.IventRepository, token string) *UserIventService {
	return &UserIventService{
		userIventRepo: userIventRepo,
		token: token,
	}
}

type UserIventCreateInput struct {
	UserIventID    uuid.UUID
	IventID        uuid.UUID
	NumberOfGuests int64
}

func (s *UserIventService) CreateUserIvent(ctx context.Context, initData string, input UserIventCreateInput) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	err = s.iventRepo.UpdateIvent(ctx, input.IventID, input.NumberOfGuests)
	if err != nil {
		return err
	}

	userIvent := &models.UserIvent{
		ID: input.UserIventID,
		UserID: user.ID,
		IventID: input.IventID,
		NumberOfGuests: input.NumberOfGuests,
	}

	err = s.userIventRepo.CreateUserIvent(ctx, userIvent)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserIventService) GetUserIventInfo(ctx context.Context, initData string, id uuid.UUID) (*models.UserIvent, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	userIvent, err := s.userIventRepo.GetUserIventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIventNotFound
		}
		return nil, err
	}
	
	return userIvent, nil 
}

func (s *UserIventService) GetUserIvents(ctx context.Context, initData string) ([]*models.UserIvent, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	userIvents, err := s.userIventRepo.GetUserIvents(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return userIvents, nil
}

func (s *UserIventService) UpdateUserIvent(ctx context.Context, initData string, id uuid.UUID, newGuestNumber int64) (*models.UserIvent, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	userIvent, err := s.userIventRepo.GetUserIventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIventNotFound
		}
		return nil, err
	}

	err = s.iventRepo.UpdateIvent(ctx, userIvent.ID, newGuestNumber-userIvent.NumberOfGuests)
	if err != nil {
		return nil, err
	}

	userIvent, err = s.userIventRepo.UpdateUserIvent(ctx, id, newGuestNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrIventNotFound
		}
		return nil, err
	}

	return userIvent, nil
}

func (s *UserIventService) DeleteUserIvent(ctx context.Context, initData string, id uuid.UUID) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	userIvent, err := s.userIventRepo.GetUserIventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrIventNotFound
		}
		return err
	}

	err = s.iventRepo.UpdateIvent(ctx, userIvent.ID, -userIvent.NumberOfGuests)
	if err != nil {
		return err
	}

	err = s.userIventRepo.DeleteUserIvent(ctx, id)
	if err != nil {
		return err
	}

	return nil
}