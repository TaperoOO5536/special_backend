package api

import (
	"context"

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
			Picture: picture.Path,
			MimeType: picture.MimeType,
		}
		pbPictures = append(pbPictures, pbPicture)
	}

	return &pb.GetEventInfoResponse{
		Id:            event.ID.String(),
		Title:         event.Title,
		Description:   event.Description,
		Datetime:      timestamppb.New(event.DateTime),
		Price:         event.Price,
		TotalSeats:    event.TotalSeats,
		OccupiedSeats: event.OccupiedSeats,
		Pictures:      pbPictures,
	}, nil
}

func (h *EventServiceHandler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.eventService.GetEvents(ctx)
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	response := &pb.GetEventsResponse{
		Events: make([]*pb.EventInfoForList, 0, len(events)),
	}
	
	for _, event := range events {
		pbEvent := &pb.EventInfoForList{
			Id:      event.ID.String(),
			Title:   event.Title,
			Price:   event.Price,
			Picture: &pb.PictureInfo{
			Picture: event.LittlePicture,
			MimeType: event.MimeType,
		},
		}
		response.Events = append(response.Events, pbEvent)
	}

	return response, nil
}