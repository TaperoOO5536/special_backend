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

type UserEventServiceHandler struct {
	userEventService *service.UserEventService
}

func NewUserEventServiceHandler(userEventService *service.UserEventService) *UserEventServiceHandler {
	return &UserEventServiceHandler{ userEventService: userEventService}
}

func (h *UserEventServiceHandler) CreateUserEvent(ctx context.Context, req *pb.CreateUserEventRequest) (*emptypb.Empty, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.EventId == "" {
		err := status.Error(codes.InvalidArgument, "event id is required")
		return nil, err
	}

	if req.NumberOfGuests == 0 {
		err := status.Error(codes.InvalidArgument, "number of guests is required")
		return nil, err
	}

	EventID, err := uuid.Parse(req.EventId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid event id")
		return nil, err
	}

	userEventID := uuid.New()

	input := service.UserEventCreateInput{
		UserEventID:    userEventID,
		EventID:        EventID,
		NumberOfGuests: int64(req.NumberOfGuests),
	}

	err = h.userEventService.CreateUserEvent(ctx, initData, input)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

func (h *UserEventServiceHandler) GetUserEventInfo(ctx context.Context, req *pb.GetUserEventInfoRequest) (*pb.GetUserEventInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user event id is required")
		return nil, err
	}	

	userEventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user event id")
		return nil, err
	}

	userEvent, err := h.userEventService.GetUserEventInfo(ctx, initData, userEventID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	return &pb.GetUserEventInfoResponse{
		Id:             userEvent.ID.String(),
		EventId:        userEvent.EventID.String(),
		NumberOfGuests: int32(userEvent.NumberOfGuests),
		Title:          userEvent.Event.Title,
		Datetime:       timestamppb.New(userEvent.Event.DateTime),
		Picture:        &pb.PictureInfo{
			Picture:  userEvent.Event.LittlePicture,
			MimeType: userEvent.Event.MimeType,
		},
	}, nil
}

func (h *UserEventServiceHandler) GetUserEvents(ctx context.Context, req *pb.GetUserEventsRequest) (*pb.GetUserEventsResponse, error) {
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

	paginatedUserEvents, err := h.userEventService.GetUserEvents(ctx, initData, pagination)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}

	pbUserEvents := make([]*pb.UserEventInfoForList, 0, len(paginatedUserEvents.UserEvents))
	for _, userEvent := range paginatedUserEvents.UserEvents {
		pbUserEvent := &pb.UserEventInfoForList{
			Id:       userEvent.ID.String(),
			EventId:  userEvent.EventID.String(),
			Title:    userEvent.Event.Title,
			Datetime: timestamppb.New(userEvent.Event.DateTime),
			Picture:  &pb.PictureInfo{
			Picture:  userEvent.Event.LittlePicture,
			MimeType: userEvent.Event.MimeType,
		},
		}
		pbUserEvents = append(pbUserEvents, pbUserEvent)
	}

	return &pb.GetUserEventsResponse{
		UserEvents: pbUserEvents,
		Total:   int32(paginatedUserEvents.TotalCount),
		Page:    int32(paginatedUserEvents.Page),
		PerPage: int32(paginatedUserEvents.PerPage),
	}, nil
}

func (h *UserEventServiceHandler) UpdateUserEvent(ctx context.Context, req *pb.UpdateUserEventRequest) (*pb.GetUserEventInfoResponse, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user event id is required")
		return nil, err
	}	

	if req.NumberOfGuests == 0 {
		err := status.Error(codes.InvalidArgument, "number of guests are required")
		return nil, err
	}

	userEventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user event id")
		return nil, err
	}

	userEvent, err := h.userEventService.UpdateUserEvent(ctx, initData, userEventID, int64(req.NumberOfGuests))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to process initData: %v", err)
	}
	
	return &pb.GetUserEventInfoResponse{
		Id:             userEvent.ID.String(),
		EventId:        userEvent.EventID.String(),
		NumberOfGuests: int32(userEvent.NumberOfGuests),
		Title:          userEvent.Event.Title,
		Datetime:       timestamppb.New(userEvent.Event.DateTime),
		Picture:        &pb.PictureInfo{
			Picture:  userEvent.Event.LittlePicture,
			MimeType: userEvent.Event.MimeType,
		},
	}, nil
}

func (h *UserEventServiceHandler) DeleteUserEvent(ctx context.Context, req *pb.DeleteUserEventRequest) (*emptypb.Empty, error) {
	initData, err := GetInitDataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "user event id is required")
		return nil, err
	}

	userEventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid user event id")
		return nil, err
	}

	err = h.userEventService.DeleteUserEvent(ctx, initData, userEventID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}