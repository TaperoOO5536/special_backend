package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/TaperoOO5536/special_backend/internal/kafka"
	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderService struct {
	orderRepo repository.OrderRepository
	token     string
	producer  *kafka.Producer
}

func NewOrderService(orderRepo repository.OrderRepository, token string, producer *kafka.Producer) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		token: token,
		producer: producer,
	}
}

type OrderCreateInput struct {
	OrderID        uuid.UUID
	CompletionDate time.Time
	Comment        string
	OrderItems     []models.OrderItem
	OrderAmount    int64
}

func (s *OrderService) CreateOrder(ctx context.Context, initData string, input OrderCreateInput) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	order := &models.Order{
		ID:             input.OrderID,
		UserID:         user.ID,
		FormDate:       time.Now(),
		CompletionDate: input.CompletionDate,
		Comment:        input.Comment,
		Status:         "В обработке",
		OrderAmount:    input.OrderAmount,
		OrderItems:     input.OrderItems,
	}

	err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	createdOrder, err := s.orderRepo.GetOrderInfo(ctx, input.OrderID)
	if err != nil {
		return err
	}

	go func() {
		msg := models.CreatedOrder{
			Number:         createdOrder.Number,
			UserID:         user.ID,
			CompletionDate: createdOrder.CompletionDate,
			OrderAmount:    createdOrder.OrderAmount,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
				return
		}

		err = s.producer.Produce(
			string(jsonMsg),
			"orders",
			"order.create",
		)
		if err != nil {
				log.Printf("failed to produce message: %v", err)
				return
		}
	}()
	
	return nil
}

func (s *OrderService) GetOrderInfo(ctx context.Context, initData string, id uuid.UUID) (*models.Order, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	order, err := s.orderRepo.GetOrderInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	
	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, initData string, pagination models.Pagination) (*models.PaginatedOrders, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.GetOrders(ctx, user.ID, pagination)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	return orders, nil
}