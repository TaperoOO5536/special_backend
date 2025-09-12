package api

import (
	"context"

	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	// "google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedSpecialAppServiceServer
	itemHandler      *ItemServiceHandler
	iventHandler     *IventServiceHandler
	userHandler      *UserServiceHandler
	orderHandler     *OrderServiceHandler
	userIventHandler *UserIventServiceHandler
}

func NewHandler(
	itemHandlel      *ItemServiceHandler,
	ivetnHandler     *IventServiceHandler,
	userHandler      *UserServiceHandler,
	orderHandler     *OrderServiceHandler,
	userIventHandler *UserIventServiceHandler,
) *Handler {
	return &Handler{
		itemHandler:      itemHandlel,
		iventHandler:     ivetnHandler,
		userHandler:      userHandler,
		orderHandler:     orderHandler,
		userIventHandler: userIventHandler,
	}
}

//items

func (h *Handler) GetItemInfo(ctx context.Context, req *pb.GetItemInfoRequest) (*pb.GetItemInfoResponse, error) {
	return h.itemHandler.GetItemInfo(ctx, req)
}

func (h *Handler) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	return h.itemHandler.GetItems(ctx, req)
}


//ivents

func (h *Handler) GetIventInfo(ctx context.Context, req *pb.GetIventInfoRequest) (*pb.GetIventInfoResponse, error) {
	return h.iventHandler.GetIventInfo(ctx, req)
}

func (h *Handler) GetIvents(ctx context.Context, req *pb.GetIventsRequest) (*pb.GetIventsResponse, error) {
	return h.iventHandler.GetIvents(ctx, req)
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

//userivents

func (h *Handler) CreateUserIvent(ctx context.Context, req *pb.CreateUserIventRequest) (*emptypb.Empty, error) {
	return h.userIventHandler.CreateUserIvent(ctx, req)
}

func (h *Handler) GetUserIventInfo(ctx context.Context, req *pb.GetUserIventInfoRequest) (*pb.GetUserIventInfoResponse, error) {
	return h.userIventHandler.GetUserIventInfo(ctx, req)
}

func (h *Handler) GetUserIvents(ctx context.Context, req *pb.GetUserIventsRequest) (*pb.GetUserIventsResponse, error) {
	return h.userIventHandler.GetUserIvents(ctx, req)
}

func (h *Handler) UpdateUserIvent(ctx context.Context, req *pb.UpdateUserIventRequest) (*pb.GetUserIventInfoResponse, error) {
	return h.userIventHandler.UpdateUserIvent(ctx, req)
}

func (h *Handler) DeleteUserIvent(ctx context.Context, req *pb.DeleteUserIventRequest) (*emptypb.Empty, error) {
	return h.userIventHandler.DeleteUserIvent(ctx, req)
}