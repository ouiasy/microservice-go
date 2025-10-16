package grpc

import (
	"context"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
	"github.com/ouiasy/microservice-go/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppState) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	err := putRating(ctx, s.db,
		model.RecordID(req.RecordId),
		model.RecordType(req.RecordType),
		&model.Rating{
			RecordID:   model.RecordID(req.RecordId),
			RecordType: model.RecordType(req.RecordType),
			UserID:     model.UserID(req.UserId),
			Value:      model.RatingValue(req.RatingValue),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PutRatingResponse{}, nil
}

func putRating(ctx context.Context, db ratingRepository, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return db.Put(ctx, recordID, recordType, rating)
}
