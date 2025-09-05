package api

import (
	"context"

	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	// "google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedSpecialAppServiceServer
	itemHandler *ItemServiceHandler
	iventHandler *IventServiceHandler
}

func NewHandler(
	itemHandlel *ItemServiceHandler,
	ivetnHandler *IventServiceHandler,
) *Handler {
	return &Handler{
		itemHandler: itemHandlel,
		iventHandler: ivetnHandler,
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