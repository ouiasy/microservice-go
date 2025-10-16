package grpc

import (
	"context"
	"errors"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
	"github.com/ouiasy/microservice-go/rating/internal/repository"
	"github.com/ouiasy/microservice-go/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppState) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required parameters")
	}
	v, err := getAggregatedRating(ctx, s.db, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

func getAggregatedRating(ctx context.Context, db ratingRepository, id model.RecordID, typ model.RecordType) (float64, error) {
	ratings, err := db.Get(ctx, model.RecordID(id), model.RecordType(typ))
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}
