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
	// UpdateOrder(ctx context.Context, id uuid.UUID, newStatus string) (*models.Order, error)
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
	if err := r.db.Preload("OrderItems.Item", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_item", "item_title", "item_price", "little_picture")
	}).Where("id_order = ?", id).First(&order).Error; err != nil {
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

// func (r *orderRepository) UpdateOrder(ctx context.Context, id uuid.UUID, newStatus string) (*models.Order, error) {
// 	order, err := r.GetOrderInfo(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	order.Status = newStatus
// 	if err := r.db.Save(order).Error; err != nil {
// 		return nil, err
// 	}
// 	return order, nil
// }