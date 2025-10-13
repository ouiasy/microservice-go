package main

import (
	"log"
	"net"

	"github.com/ouiasy/microservice-go/metadata/internal/controller"
	grpcHandler "github.com/ouiasy/microservice-go/metadata/internal/handler/grpc"
	"github.com/ouiasy/microservice-go/metadata/internal/repository/memory"
	gen "github.com/ouiasy/microservice/_proto/gen/go"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := controller.New(repo)
	h := grpcHandler.New(ctrl)
	server := grpc.NewServer()
	gen.RegisterMetadataServiceServer(server, h)

	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
