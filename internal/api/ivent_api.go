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

type IventServiceHandler struct {
	iventService *service.IventService
}

func NewIventServiceHandler(iventService *service.IventService) *IventServiceHandler {
	return &IventServiceHandler{ iventService: iventService }
}

func (h *IventServiceHandler) GetIventInfo(ctx context.Context, req *pb.GetIventInfoRequest) (*pb.GetIventInfoResponse, error) {
	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "ivent ID is required")
		return nil, err
	}
	
	iventID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid ivent ID")
		return nil, err
	}

	ivent, err := h.iventService.GetIventInfo(ctx, iventID)
	if err != nil {
		if err == service.ErrIventNotFound {
			err := status.Error(codes.NotFound, "ivent not found")
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	var pictures []string
	for _, picture := range ivent.Pictures {
		pictures = append(pictures, picture.Path)
	}

	return &pb.GetIventInfoResponse{
		Id:            ivent.ID.String(),
		Title:         ivent.Title,
		Description:   ivent.Description,
		Datetime:      timestamppb.New(ivent.DateTime),
		Price:         ivent.Price,
		TotalSeats:    ivent.TotalSeats,
		OccupiedSeats: ivent.OccupiedSeats,
		Pictures:      pictures,
	}, nil
}

func (h *IventServiceHandler) GetIvents(ctx context.Context, req *pb.GetIventsRequest) (*pb.GetIventsResponse, error) {
	ivents, err := h.iventService.GetIvents(ctx)
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	response := &pb.GetIventsResponse{
		Ivents: make([]*pb.IventInfoForList, 0, len(ivents)),
	}
	
	for _, ivent := range ivents {
		pbIvent := &pb.IventInfoForList{
			Id:      ivent.ID.String(),
			Title:   ivent.Title,
			Price:   ivent.Price,
			Picture: ivent.LittlePicture,
		}
		response.Ivents = append(response.Ivents, pbIvent)
	}

	return response, nil
}