package api

import (
	"context"

	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	// "google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedSpecialAppServiceServer
	itemHandler *ItemServiceHandler
}

func NewHandler(
	itemHandlel *ItemServiceHandler,
) *Handler {
	return &Handler{
		itemHandler: itemHandlel,
	}
}

func (h *Handler) GetItemInfo(ctx context.Context, req *pb.GetItemInfoRequest) (*pb.GetItemInfoResponse, error) {
	return h.itemHandler.GetItemInfo(ctx, req)
}