package api

import (
	"context"

	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedSpecialAppServiceServer
	itemHandler      *ItemServiceHandler
	eventHandler     *EventServiceHandler
	userHandler      *UserServiceHandler
	orderHandler     *OrderServiceHandler
	userEventHandler *UserEventServiceHandler
}

func NewHandler(
	itemHandlel      *ItemServiceHandler,
	ivetnHandler     *EventServiceHandler,
	userHandler      *UserServiceHandler,
	orderHandler     *OrderServiceHandler,
	userEventHandler *UserEventServiceHandler,
) *Handler {
	return &Handler{
		itemHandler:      itemHandlel,
		eventHandler:     ivetnHandler,
		userHandler:      userHandler,
		orderHandler:     orderHandler,
		userEventHandler: userEventHandler,
	}
}

//items

func (h *Handler) GetItemInfo(ctx context.Context, req *pb.GetItemInfoRequest) (*pb.GetItemInfoResponse, error) {
	return h.itemHandler.GetItemInfo(ctx, req)
}

func (h *Handler) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	return h.itemHandler.GetItems(ctx, req)
}


//events

func (h *Handler) GetEventInfo(ctx context.Context, req *pb.GetEventInfoRequest) (*pb.GetEventInfoResponse, error) {
	return h.eventHandler.GetEventInfo(ctx, req)
}

func (h *Handler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return h.eventHandler.GetEvents(ctx, req)
}

//users

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserInfoResponse, error) {
	return h.userHandler.GetUserInfo(ctx, req)
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	return h.userHandler.CreateUser(ctx, req)
}

func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserInfoResponse, error) {
	return h.userHandler.UpdateUser(ctx, req)
}

//orders

func (h *Handler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*emptypb.Empty, error) {
	return h.orderHandler.CreateOrder(ctx, req)
}

func (h *Handler) GetOrderInfo(ctx context.Context, req *pb.GetOrderInfoRequest) (*pb.GetOrderInfoResponse, error) {
	return h.orderHandler.GetOrderInfo(ctx, req)
}

func (h *Handler) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return h.orderHandler.GetOrders(ctx, req)
}

//userevents

func (h *Handler) CreateUserEvent(ctx context.Context, req *pb.CreateUserEventRequest) (*emptypb.Empty, error) {
	return h.userEventHandler.CreateUserEvent(ctx, req)
}

func (h *Handler) GetUserEventInfo(ctx context.Context, req *pb.GetUserEventInfoRequest) (*pb.GetUserEventInfoResponse, error) {
	return h.userEventHandler.GetUserEventInfo(ctx, req)
}

func (h *Handler) GetUserEvents(ctx context.Context, req *pb.GetUserEventsRequest) (*pb.GetUserEventsResponse, error) {
	return h.userEventHandler.GetUserEvents(ctx, req)
}

func (h *Handler) UpdateUserEvent(ctx context.Context, req *pb.UpdateUserEventRequest) (*pb.GetUserEventInfoResponse, error) {
	return h.userEventHandler.UpdateUserEvent(ctx, req)
}

func (h *Handler) DeleteUserEvent(ctx context.Context, req *pb.DeleteUserEventRequest) (*emptypb.Empty, error) {
	return h.userEventHandler.DeleteUserEvent(ctx, req)
}