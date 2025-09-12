package api

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"google.golang.org/grpc/codes"
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
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.userService.CreateUser(ctx, initData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *UserServiceHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

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

func (h *UserServiceHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var phoneNumber string
	if req.PhoneNumber != nil {
		phoneNumber = *req.PhoneNumber
	} else {
		phoneNumber = ""
	}

	user, err := h.userService.UpdateUser(ctx, initData, phoneNumber)
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