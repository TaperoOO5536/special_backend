package api

import (
	"context"

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
		return nil, err
	}

	if req.CompletionDate == nil {
		err := status.Error(codes.InvalidArgument, "completion date is required")
		return nil, err
	}

	if req.Items == nil || len(req.Items) == 0 {
		err := status.Error(codes.InvalidArgument, "order items are required")
		return nil, err
	}

	if req.OrderAmount == 0 {
		err := status.Error(codes.InvalidArgument, "order amount is required")
		return nil, err
	}

	orderID := uuid.New()

	var items []models.OrderItem
	for _, pbItem := range req.Items {
		itemID, err := uuid.Parse(pbItem.ItemId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid itemID")
		}
		item := models.OrderItem{
			ID: uuid.New(),
			OrderID: orderID,
			ItemID: itemID,
			Quantity: pbItem.Quantity,
		}
		items = append(items, item)
	}

	input := service.OrderCreateInput{
		OrderID: orderID,
		CompletionDate: req.CompletionDate.AsTime(),
		Comment: *req.Comment,
		OrderItems: items,
		OrderAmount: req.OrderAmount,
	}

	err = h.orderService.CreateOrder(ctx, initData, input)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *OrderServiceHandler) GetOrderInfo(ctx context.Context, req *pb.GetOrderInfoRequest) (*pb.GetOrderInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "order id is required")
		return nil, err
	}	

	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid order id")
		return nil, err
	}

	order, err := h.orderService.GetOrderInfo(ctx, initData, orderID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	orderItems := make([]*pb.OrderItemInfoForList, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		pbItem := &pb.OrderItemInfoForList{
			Id: orderItem.ID.String(),
			ItemId: orderItem.ItemID.String(),
			Title: orderItem.Item.Title,
			Price: orderItem.Item.Price,
			Quantity: orderItem.Quantity,
			Picture: &pb.PictureInfo{
			Picture: orderItem.Item.LittlePicture,
			MimeType: orderItem.Item.MimeType,
		},
		}
		orderItems = append(orderItems, pbItem)
	}
	
	return &pb.GetOrderInfoResponse{
		Number: order.Number,
		FormDate: timestamppb.New(order.FormDate),
		CompletionDate: timestamppb.New(order.CompletionDate),
		Comment: order.Comment,
		Status: order.Status,
		OrderAmount: order.OrderAmount,
		Items: orderItems,
	}, nil
}

func (h *OrderServiceHandler) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
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
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	pbOrders := make([]*pb.OrderInfoForList, 0, len(paginatedOrders.Orders))
	for _, order := range paginatedOrders.Orders {
		pbOrder := &pb.OrderInfoForList{
			Number: order.Number,
			CompletionDate: timestamppb.New(order.CompletionDate),
			Status: order.Status,
			OrderAmount: order.OrderAmount,
		}
		pbOrders = append(pbOrders, pbOrder)
	}

	return &pb.GetOrdersResponse{
		Orders: pbOrders,
		Total:   paginatedOrders.TotalCount,
		Page:    int32(paginatedOrders.Page),
		PerPage: int32(paginatedOrders.PerPage),
	}, nil
}