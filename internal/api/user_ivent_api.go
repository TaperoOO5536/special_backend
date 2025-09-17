package api

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserIventServiceHandler struct {
	userIventService *service.UserIventService
}

func NewUserIventServiceHandler(userIventService *service.UserIventService) *UserIventServiceHandler {
	return &UserIventServiceHandler{ userIventService: userIventService}
}

func (h *UserIventServiceHandler) CreateUserIvent(ctx context.Context, req *pb.CreateUserIventRequest) (*emptypb.Empty, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.IventId == "" {
		err := status.Error(codes.InvalidArgument, "ivent id is required")
		return nil, err
	}

	if req.NumberOfGuests == 0 {
		err := status.Error(codes.InvalidArgument, "number of guests are required")
		return nil, err
	}

	IventID, err := uuid.Parse(req.IventId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid ivent id")
		return nil, err
	}

	userIventID := uuid.New()

	input := service.UserIventCreateInput{
		UserIventID: userIventID,
		IventID: IventID,
		NumberOfGuests: req.NumberOfGuests,
	}

	err = h.userIventService.CreateUserIvent(ctx, initData, input)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *UserIventServiceHandler) GetUserIventInfo(ctx context.Context, req *pb.GetUserIventInfoRequest) (*pb.GetUserIventInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user ivent id is required")
		return nil, err
	}	

	userIventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user ivent id")
		return nil, err
	}

	userIvent, err := h.userIventService.GetUserIventInfo(ctx, initData, userIventID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	return &pb.GetUserIventInfoResponse{
		Id: userIvent.ID.String(),
		IventId: userIvent.IventID.String(),
		NumberOfGuests: userIvent.NumberOfGuests,
		Title: userIvent.Ivent.Title,
		Datetime: timestamppb.New(userIvent.Ivent.DateTime),
		Picture: &pb.PictureInfo{
			Picture: userIvent.Ivent.LittlePicture,
			MimeType: userIvent.Ivent.MimeType,
		},
	}, nil
}

func (h *UserIventServiceHandler) GetUserIvents(ctx context.Context, req *pb.GetUserIventsRequest) (*pb.GetUserIventsResponse, error) {
 initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userIvents, err := h.userIventService.GetUserIvents(ctx, initData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	pbUserIvents := make([]*pb.UserIventInfoForList, 0, len(userIvents))
	for _, userIvent := range userIvents {
		pbUserIvent := &pb.UserIventInfoForList{
			Id: userIvent.ID.String(),
			IventId: userIvent.IventID.String(),
			Title: userIvent.Ivent.Title,
			Datetime: timestamppb.New(userIvent.Ivent.DateTime),
			Picture: &pb.PictureInfo{
			Picture: userIvent.Ivent.LittlePicture,
			MimeType: userIvent.Ivent.MimeType,
		},
		}
		pbUserIvents = append(pbUserIvents, pbUserIvent)
	}

	return &pb.GetUserIventsResponse{
		UserIvents: pbUserIvents,
	}, nil
}

func (h *UserIventServiceHandler) UpdateUserIvent(ctx context.Context, req *pb.UpdateUserIventRequest) (*pb.GetUserIventInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user ivent id is required")
		return nil, err
	}	

	if req.NumberOfGuests == 0 {
		err := status.Error(codes.InvalidArgument, "number of guests are required")
		return nil, err
	}

	userIventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user ivent id")
		return nil, err
	}

	userIvent, err := h.userIventService.UpdateUserIvent(ctx, initData, userIventID, req.NumberOfGuests)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &pb.GetUserIventInfoResponse{
		Id: userIvent.ID.String(),
		IventId: userIvent.IventID.String(),
		NumberOfGuests: userIvent.NumberOfGuests,
		Title: userIvent.Ivent.Title,
		Datetime: timestamppb.New(userIvent.Ivent.DateTime),
		Picture: &pb.PictureInfo{
			Picture: userIvent.Ivent.LittlePicture,
			MimeType: userIvent.Ivent.MimeType,
		},
	}, nil
}

func (h *UserIventServiceHandler) DeleteUserIvent(ctx context.Context, req *pb.DeleteUserIventRequest) (*emptypb.Empty, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user ivent id is required")
		return nil, err
	}

	userIventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user ivent id")
		return nil, err
	}

	err = h.userIventService.DeleteUserIvent(ctx, initData, userIventID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}