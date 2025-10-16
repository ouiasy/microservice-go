package main

import (
	"log"
	"net"

	gen "github.com/ouiasy/microservice-go/common/gen/go"
	grpcAppState "github.com/ouiasy/microservice-go/rating/internal/handler/grpc"
	"github.com/ouiasy/microservice-go/rating/internal/repository/memory"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()

	s := grpcAppState.NewAppState(repo)
	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	gen.RegisterRatingServiceServer(server, s)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
