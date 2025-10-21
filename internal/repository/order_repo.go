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
	GetOrders(ctx context.Context, userID string, pagination models.Pagination) (*models.PaginatedOrders, error)
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
		return db.Select("id_item", "item_title", "item_price", "little_picture", "mime_type")
	}).Where("id_order = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrders(ctx context.Context, userID string, pagination models.Pagination) (*models.PaginatedOrders, error) {
	var orders []models.Order
	var total int64
	
	if err := r.db.Model(&models.Order{}).Count(&total).Error; err != nil {
    return nil, err
  }

	offset := (pagination.Page - 1) * pagination.PerPage
	if err := r.db.Where("user_id = ?", userID).Limit(pagination.PerPage).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &models.PaginatedOrders{
		Orders:     orders,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	}, nil
}