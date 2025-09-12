package repository

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderInfo(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetOrders(ctx context.Context, userID string) ([]*models.Order, error)
	// UpdateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) (error) {
	if err := r.db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) GetOrderInfo(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("OrderItems").Where("id_order = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}