package service

import (
	"context"
	"errors"
	"time"

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
	token    string
}

func NewOrderService(orderRepo repository.OrderRepository, token string) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		token: token,
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
		ID: input.OrderID,
		Number: "number",
		UserID: user.ID,
		FormDate: time.Now(),
		CompletionDate: input.CompletionDate,
		Comment: input.Comment,
		Status: "В обработке",
		OrderAmount: input.OrderAmount,
		OrderItems: input.OrderItems,
	}

	err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}
	
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

func (s *OrderService) GetOrders(ctx context.Context, initData string) ([]*models.Order, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, err
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.GetOrders(ctx, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	return orders, nil
}

// func (s *OrderService) UpdateOrder(ctx context.Context, initData string, id uuid.UUID, newStatus string) (*models.Order, error) {
// 	valid, err := VerifyInitData(initData, s.token)
// 	if err != nil || !valid {
// 		return nil, err
// 	}

// 	order, err := s.orderRepo.UpdateOrder(ctx, id, newStatus)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, ErrEventNotFound
// 		}
// 		return nil, err
// 	}

// 	return order, nil
// }