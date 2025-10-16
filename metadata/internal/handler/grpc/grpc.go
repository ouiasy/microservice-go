package grpcHandler

import (
	"context"
	"errors"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
	"github.com/ouiasy/microservice-go/metadata/internal/controller"
	"github.com/ouiasy/microservice-go/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{svc: ctrl}
}

func (h *Handler) GetMetadata(ctx context.Context, r *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if r == nil || r.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "metadata service is nil")
	}
	m, err := h.svc.Get(ctx, r.MovieId)
	if err != nil {
		if errors.Is(err, controller.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}
