package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TaperoOO5536/special_backend/internal/bot_handler"
	"github.com/TaperoOO5536/special_backend/internal/config"
	"github.com/TaperoOO5536/special_backend/internal/kafka"
	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"github.com/TaperoOO5536/special_backend/pkg/env"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderService struct {
	orderRepo  repository.OrderRepository
	token      string
	producer   *kafka.Producer
	yooclient  *config.Client
	botHandler *bot_handler.BotHandler
}

func NewOrderService(orderRepo repository.OrderRepository, token string, producer *kafka.Producer, yooclient *config.Client, botHandler *bot_handler.BotHandler) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		token: token,
		producer: producer,
		yooclient: yooclient,
		botHandler: botHandler,
	}
}

type OrderCreateInput struct {
	OrderID        uuid.UUID
	CompletionDate time.Time
	Comment        string
	OrderItems     []models.OrderItem
	OrderAmount    int64
}

type PaymentResponse struct {
	Id  string
	Url string
}

func (s *OrderService) CreatePayment1(ctx context.Context, amount int64, orderID string) (*PaymentResponse, error) {
    client := yookassa.NewClient(env.GetShopId(), env.GetYookassaSecret())
		paymentHandler := yookassa.NewPaymentHandler(client)
		payment, err := paymentHandler.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value: strconv.FormatInt(amount*100, 10),
				Currency: "RUB",
			},
			PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
			Confirmation: yoopayment.Redirect{
				Type:      "redirect",
				ReturnURL: "https://example.com",
			},
			Description: "Заказ выпечки #" + orderID,
			Metadata: map[string]string{"order_id": orderID},
		})
    if err != nil {
			return nil, err
		}

		if payment.Confirmation != nil {
    	confBytes, err := json.Marshal(payment.Confirmation)
    	if err != nil {
    	    return nil, fmt.Errorf("failed to marshal confirmation: %w", err)
    	}

    	var confMap map[string]interface{}
    	if err := json.Unmarshal(confBytes, &confMap); err != nil {
    	    return nil, fmt.Errorf("failed to unmarshal confirmation: %w", err)
    	}
			return &PaymentResponse{
				Id:  payment.ID,
				Url: confMap["confirmation_url"].(string),
			}, nil

		}

    return nil, fmt.Errorf("something went wrong")
}

func (s *OrderService) CreateOrder(ctx context.Context, initData string, input OrderCreateInput) error {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return fmt.Errorf("failed to verify init data %v", err)
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return fmt.Errorf("failed to parse init data %v", err)
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

	userID, _ := strconv.Atoi(user.ID)
	s.botHandler.HandleOrderMessage("order.create", userID, *createdOrder)


	// go func() {
	// 	msg := models.KafkaOrder{
	// 		Number:         strconv.FormatInt(int64(createdOrder.Number), 10),
	// 		UserID:         user.ID,
	// 		CompletionDate: createdOrder.CompletionDate,
	// 		OrderAmount:    createdOrder.OrderAmount,
	// 	}
	// 	jsonMsg, err := json.Marshal(msg)
	// 	if err != nil {
	// 		log.Printf("failed to marshal message: %v", err)
	// 		return
	// 	}

	// 	err = s.producer.Produce(
	// 		string(jsonMsg),
	// 		"orders",
	// 		"order.create",
	// 	)
	// 	if err != nil {
	// 			log.Printf("failed to produce message: %v", err)
	// 			return
	// 	}
	// }()
	
	return nil
}

func (s *OrderService) GetOrderInfo(ctx context.Context, initData string, id uuid.UUID) (*models.Order, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, fmt.Errorf("failed to verify init data %v", err)
	}

	order, err := s.orderRepo.GetOrderInfo(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	
	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, initData string, pagination models.Pagination) (*models.PaginatedOrders, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, fmt.Errorf("failed to verify init data %v", err)
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init data %v", err)
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