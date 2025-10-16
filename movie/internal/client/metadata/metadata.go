package metadata

import (
	"context"

	"github.com/ouiasy/microservice-go/common/discovery"
	gen "github.com/ouiasy/microservice-go/common/gen/go"
	"github.com/ouiasy/microservice-go/metadata/pkg/model"
	"github.com/ouiasy/microservice-go/movie/internal/grpc_util"
)

// Client defines a movie metadata gRPC gateway.
type Client struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service.
func New(registry discovery.Registry) *Client {
	return &Client{registry}
}

// Get returns movie metadata by a movie id.
func (g *Client) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpc_util.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}
