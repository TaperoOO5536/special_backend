package api

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServiceHandler struct {
	eventService *service.EventService
}

func NewEventServiceHandler(eventService *service.EventService) *EventServiceHandler {
	return &EventServiceHandler{ eventService: eventService }
}

func (h *EventServiceHandler) GetEventInfo(ctx context.Context, req *pb.GetEventInfoRequest) (*pb.GetEventInfoResponse, error) {
	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "event ID is required")
		return nil, err
	}
	
	eventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid event ID")
		return nil, err
	}

	event, err := h.eventService.GetEventInfo(ctx, eventID)
	if err != nil {
		if err == service.ErrEventNotFound {
			err := status.Error(codes.NotFound, "event not found")
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	pbPictures := make([]*pb.PictureInfo, 0, len(event.Pictures))
	for _, picture := range event.Pictures {
		pbPicture := &pb.PictureInfo{
			Picture:  picture.Path,
			MimeType: picture.MimeType,
		}
		pbPictures = append(pbPictures, pbPicture)
	}

	return &pb.GetEventInfoResponse{
		Id:            event.ID.String(),
		Title:         event.Title,
		Description:   event.Description,
		Datetime:      timestamppb.New(event.DateTime),
		Price:         int32(event.Price),
		TotalSeats:    int32(event.TotalSeats),
		OccupiedSeats: int32(event.OccupiedSeats),
		Pictures:      pbPictures,
	}, nil
}

func (h *EventServiceHandler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
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

	paginatedEvents, err := h.eventService.GetEvents(ctx, pagination)
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	response := &pb.GetEventsResponse{
		Events:  make([]*pb.EventInfoForList, 0, len(paginatedEvents.Events)),
		Total:   int32(paginatedEvents.TotalCount),
		Page:    int32(paginatedEvents.Page),
		PerPage: int32(paginatedEvents.PerPage),
	}
	
	for _, event := range paginatedEvents.Events {
		pbEvent := &pb.EventInfoForList{
			Id:            event.ID.String(),
			Title:         event.Title,
			Datetime:      timestamppb.New(event.DateTime),
			Price:         int32(event.Price),
			TotalSeats:    int32(event.TotalSeats),
			OccupiedSeats: int32(event.OccupiedSeats),
			Picture:       &pb.PictureInfo{
			Picture:  event.LittlePicture,
			MimeType: event.MimeType,
		},
		}
		response.Events = append(response.Events, pbEvent)
	}

	return response, nil
}