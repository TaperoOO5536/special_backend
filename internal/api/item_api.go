package api

import (
	"context"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/service"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ItemServiceHandler struct {
	itemService *service.ItemService
}

func NewItemServiceHandler(itemService *service.ItemService) *ItemServiceHandler {
	return &ItemServiceHandler{ itemService: itemService }
}

func (h *ItemServiceHandler) GetItemInfo(ctx context.Context, req *pb.GetItemInfoRequest) (*pb.GetItemInfoResponse, error) {
	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "item ID is required")
		return nil, err
	}
	
	itemID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid item ID")
		return nil, err
	}

	item, err := h.itemService.GetItemInfo(ctx, itemID)
	if err != nil {
		if err == service.ErrItemNotFound {
			err := status.Error(codes.NotFound, "item not found")
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	pbPictures := make([]*pb.PictureInfo, 0, len(item.Pictures))
	for _, picture := range item.Pictures {
		pbPicture := &pb.PictureInfo{
			Picture:  picture.Path,
			MimeType: picture.MimeType,
		}
		pbPictures = append(pbPictures, pbPicture)
	}

	return &pb.GetItemInfoResponse{
		Id:          item.ID.String(),
		Title:       item.Title,
		Description: item.Description,
		Price:       int32(item.Price),
		Pictures:    pbPictures,
	}, nil
}

func (h *ItemServiceHandler) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
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
	
	paginatedItems, err := h.itemService.GetItems(ctx, pagination)
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}

	response := &pb.GetItemsResponse{
		Items:   make([]*pb.ItemInfoForList, 0, len(paginatedItems.Items)),
		Total:   int32(paginatedItems.TotalCount),
		Page:    int32(paginatedItems.Page),
		PerPage: int32(paginatedItems.PerPage),
	}
	
	for _, item := range paginatedItems.Items {
		pbItem := &pb.ItemInfoForList{
			Id:       item.ID.String(),
			Title:    item.Title,
			Price:    int32(item.Price),
			Picture:  &pb.PictureInfo{
			Picture:  item.LittlePicture,
			MimeType: item.MimeType,
		},
		}
		response.Items = append(response.Items, pbItem)
	}

	return response, nil
}