package grpc

import (
	"errors"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
)

var (
	// ErrNotFound is returned when no ratings are found for a
	// record.
	ErrNotFound = errors.New("ratings not found for a record")
)

type AppState struct {
	gen.UnimplementedRatingServiceServer
	db ratingRepository
}

func NewAppState(s ratingRepository) *AppState {
	return &AppState{db: s}
}
