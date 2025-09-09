package api

import (
	"context"
	"strings"

	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceHandler struct {
	userService *service.UserService
}

func NewUserServiceHandler(userService *service.UserService) *UserServiceHandler {
	return &UserServiceHandler{ userService: userService}
}

func (h *UserServiceHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization")
	}
	initData := strings.TrimPrefix(authHeaders[0], "Bearer ")

	err := h.userService.CreateUser(ctx, initData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *UserServiceHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserInfoResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization")
	}
	initData := strings.TrimPrefix(authHeaders[0], "Bearer ")

	user, err := h.userService.GetUserInfo(ctx, initData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	return &pb.GetUserInfoResponse{
		Name:        user.Name,
		Surname:     user.Surname,
		Nickname:    user.Nickname,
		PhoneNumber: user.PhoneNumber,
	}, nil
}