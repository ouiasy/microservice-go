package grpc

import (
	"context"
	"errors"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
	"github.com/ouiasy/microservice-go/metadata/pkg/model"
	"github.com/ouiasy/microservice-go/movie/internal/controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie gRPC handler.
type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *controller.Controller
}

// New creates a new movie gRPC handler.
func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMovieDetails returns movie details by id.
func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: model.MetadataToProto(&m.Metadata),
			Rating:   *m.Rating,
		},
	}, nil
}
