package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/TaperoOO5536/special_backend/internal/kafka"
	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserEventNotFound = errors.New("userevent not found")
)

type UserEventService struct {
	eventRepo     repository.EventRepository
	userEventRepo repository.UserEventRepository
	token         string
	producer      *kafka.Producer
}

func NewUserEventService(userEventRepo repository.UserEventRepository,
												 eventRepo repository.EventRepository,
												 token string,
												 producer *kafka.Producer) *UserEventService {
	return &UserEventService{
		userEventRepo: userEventRepo,
		eventRepo: eventRepo,
		token: token,
		producer: producer,
	}
}

type UserEventCreateInput struct {
	UserEventID    uuid.UUID
	EventID        uuid.UUID
	NumberOfGuests int64
}

func (s *UserEventService) CreateUserEvent(ctx context.Context, initData string, input UserEventCreateInput) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	err = s.eventRepo.UpdateEvent(ctx, input.EventID, input.NumberOfGuests)
	if err != nil {
		return err
	}

	userEvent := &models.UserEvent{
		ID: input.UserEventID,
		UserID: user.ID,
		EventID: input.EventID,
		NumberOfGuests: input.NumberOfGuests,
	}

	err = s.userEventRepo.CreateUserEvent(ctx, userEvent)
	if err != nil {
		return err
	}

	event, err := s.eventRepo.GetEventInfo(ctx, input.EventID)
	if err != nil {
		return err
	}

	go func() {
		msg := models.KafkaUserEvent{
			UserNickName:       user.Nickname,
			EventTitle:          event.Title,
			EventOccupiedSeats: event.OccupiedSeats,
			EventTotalSeats:    event.TotalSeats,
			NumberOfGuests:     input.NumberOfGuests,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			return
		}
		err = s.producer.Produce(
			string(jsonMsg),
			"userevents",
			"userevent.create",
		)
		if err != nil {
			log.Printf("failed to produce message: %v", err)
			return
		}
	}()

	return nil
}

func (s *UserEventService) GetUserEventInfo(ctx context.Context, initData string, id uuid.UUID) (*models.UserEvent, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	userEvent, err := s.userEventRepo.GetUserEventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserEventNotFound
		}
		return nil, err
	}
	
	return userEvent, nil 
}

func (s *UserEventService) GetUserEvents(ctx context.Context, initData string, pagination models.Pagination) (*models.PaginatedUserEvents, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	userEvents, err := s.userEventRepo.GetUserEvents(ctx, user.ID, pagination)
	if err != nil {
		return nil, err
	}

	return userEvents, nil
}

func (s *UserEventService) UpdateUserEvent(ctx context.Context, initData string, id uuid.UUID, newGuestNumber int64) (*models.UserEvent, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	userEvent, err := s.userEventRepo.GetUserEventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserEventNotFound
		}
		return nil, err
	}

	err = s.eventRepo.UpdateEvent(ctx, userEvent.EventID, newGuestNumber-userEvent.NumberOfGuests)
	if err != nil {
		return nil, err
	}

	userEvent, err = s.userEventRepo.UpdateUserEvent(ctx, id, newGuestNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	event, err := s.eventRepo.GetEventInfo(ctx, userEvent.EventID)
	if err != nil {
		return nil, err
	}

	go func() {
		msg := models.KafkaUserEvent{
			UserNickName:       user.Nickname,
			EventTitle:         event.Title,
			EventOccupiedSeats: event.OccupiedSeats,
			EventTotalSeats:    event.TotalSeats,
			NumberOfGuests:     userEvent.NumberOfGuests,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			return
		}
		err = s.producer.Produce(
			string(jsonMsg),
			"userevents",
			"userevent.update",
		)
		if err != nil {
			log.Printf("failed to produce message: %v", err)
			return
		}
	}()

	return userEvent, nil
}

func (s *UserEventService) DeleteUserEvent(ctx context.Context, initData string, id uuid.UUID) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	userEvent, err := s.userEventRepo.GetUserEventInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrEventNotFound
		}
		return err
	}

	err = s.eventRepo.UpdateEvent(ctx, userEvent.EventID, -userEvent.NumberOfGuests)
	if err != nil {
		return err
	}

	err = s.userEventRepo.DeleteUserEvent(ctx, id)
	if err != nil {
		return err
	}

	event, err := s.eventRepo.GetEventInfo(ctx, userEvent.EventID)
	if err != nil {
		return err
	}

	go func() {
		msg := models.KafkaUserEvent{
			UserNickName:       user.Nickname,
			EventTitle:         event.Title,
			EventOccupiedSeats: event.OccupiedSeats,
			EventTotalSeats:    event.TotalSeats,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			return
		}
		err = s.producer.Produce(
			string(jsonMsg),
			"userevents",
			"userevent.delete",
		)
		if err != nil {
			log.Printf("failed to produce message: %v", err)
			return
		}
	}()

	return nil
}