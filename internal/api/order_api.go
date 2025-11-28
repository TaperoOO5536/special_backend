package api

import (
	"context"
	"log"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServiceHandler struct {
	orderService *service.OrderService
}

func NewOrderServiceHandler(orderService *service.OrderService) *OrderServiceHandler {
	return &OrderServiceHandler{ orderService: orderService}
}

func (h *OrderServiceHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*emptypb.Empty, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.CompletionDate == nil {
		err := status.Error(codes.InvalidArgument, "completion date is required")
		log.Println(err)
		return nil, err
	}

	if req.Items == nil || len(req.Items) == 0 {
		err := status.Error(codes.InvalidArgument, "order items are required")
		log.Println(err)
		return nil, err
	}

	if req.OrderAmount == 0 {
		err := status.Error(codes.InvalidArgument, "order amount is required")
		log.Println(err)
		return nil, err
	}

	orderID := uuid.New()

	var items []models.OrderItem
	for _, pbItem := range req.Items {
		itemID, err := uuid.Parse(pbItem.ItemId)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.InvalidArgument, "invalid itemID")
		}
		item := models.OrderItem{
			ID:       uuid.New(),
			OrderID:  orderID,
			ItemID:   itemID,
			Quantity: int64(pbItem.Quantity),
		}
		items = append(items, item)
	}

	input := service.OrderCreateInput{
		OrderID:        orderID,
		CompletionDate: req.CompletionDate.AsTime(),
		Comment:        *req.Comment,
		OrderItems:     items,
		OrderAmount:    int64(req.OrderAmount),
	}

	err = h.orderService.CreateOrder(ctx, initData, input)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *OrderServiceHandler) GetOrderInfo(ctx context.Context, req *pb.GetOrderInfoRequest) (*pb.GetOrderInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "order id is required")
		log.Println(err)
		return nil, err
	}	

	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid order id")
		log.Println(err)
		return nil, err
	}

	order, err := h.orderService.GetOrderInfo(ctx, initData, orderID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	orderItems := make([]*pb.OrderItemInfoForList, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		pbItem := &pb.OrderItemInfoForList{
			Id:       orderItem.ID.String(),
			ItemId:   orderItem.ItemID.String(),
			Title:    orderItem.Item.Title,
			Price:    int32(orderItem.Item.Price),
			Quantity: int32(orderItem.Quantity),
			Picture:  orderItem.Item.LittlePicture,
		}
		orderItems = append(orderItems, pbItem)
	}
	
	return &pb.GetOrderInfoResponse{
		Number:         order.Number,
		FormDate:       timestamppb.New(order.FormDate),
		CompletionDate: timestamppb.New(order.CompletionDate),
		Comment:        order.Comment,
		Status:         order.Status,
		OrderAmount:    int32(order.OrderAmount),
		Items:          orderItems,
	}, nil
}

func (h *OrderServiceHandler) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pagination := models.Pagination{}	

	if req.Page == 0 {
		pagination.Page = 1
	} else {
		pagination.Page = int(req.Page)
	}
	if req.PerPage == 0 {
		pagination.PerPage = 1
	} else {
		pagination.PerPage = int(req.PerPage)
	}

	paginatedOrders, err := h.orderService.GetOrders(ctx, initData, pagination)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}
	if paginatedOrders == nil {
		err := status.Error(codes.Internal, "internal error: no orders returned")
		log.Println(err)
    return nil, err
  }
  if paginatedOrders.Orders == nil {
    paginatedOrders.Orders = []models.Order{}
  }

	pbOrders := make([]*pb.OrderInfoForList, 0, len(paginatedOrders.Orders))
	for _, order := range paginatedOrders.Orders {
		pbOrder := &pb.OrderInfoForList{
			Id:             order.ID.String(),
			Number:         order.Number,
			CompletionDate: timestamppb.New(order.CompletionDate),
			Status:         order.Status,
			OrderAmount:    int32(order.OrderAmount),
		}
		pbOrders = append(pbOrders, pbOrder)
	}

	return &pb.GetOrdersResponse{
		Orders: pbOrders,
		Total:   int32(paginatedOrders.TotalCount),
		Page:    int32(paginatedOrders.Page),
		PerPage: int32(paginatedOrders.PerPage),
	}, nil
}

func (h *OrderServiceHandler) Payment(ctx context.Context, req *pb.GetPaymentUrlRequest) (*pb.GetPaymentUrlResponse, error) {
	// initData, err := GetInitDataFromContext(ctx)
	// if err != nil {
	// 	return nil
	// }

	if req.OrderId == "" {
		err := status.Error(codes.InvalidArgument, "order id is required")
		return nil, err
	}	

	if req.Amount == 0 {
		err := status.Error(codes.InvalidArgument, "order amount is required")
		return nil, err
	}	

	resp, err := h.orderService.CreatePayment1(ctx, req.Amount, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &pb.GetPaymentUrlResponse{
		PaymentId: resp.Id,
		PaymentUrl: resp.Url,
	}, nil
}