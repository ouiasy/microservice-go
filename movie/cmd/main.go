package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ouiasy/microservice-go/common/discovery"
	gen "github.com/ouiasy/microservice-go/common/gen/go"
	metadataClient "github.com/ouiasy/microservice-go/movie/internal/client/metadata"
	ratingClient "github.com/ouiasy/microservice-go/movie/internal/client/rating"
	"github.com/ouiasy/microservice-go/movie/internal/controller"
	grpc_handler "github.com/ouiasy/microservice-go/movie/internal/handler/grpc"
	"google.golang.org/grpc"
)

const serviceName = "movie"

func main() {
	log.Println("Starting the movie service")

	consulClient, err := discovery.NewConsulClient("consul:8500")
	if err != nil {
		log.Fatal(err)
	}
	instanceId := discovery.GenerateInstanceID(serviceName)
	ctx := context.Background()
	if err := consulClient.Register(ctx, instanceId, serviceName, "localhost:8083"); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			if err := consulClient.ReportHealthyState(instanceId, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer consulClient.Deregister(ctx, instanceId, serviceName)

	metadataGateway := metadataClient.New(consulClient)
	ratingGateway := ratingClient.New(consulClient)
	ctrl := controller.New(ratingGateway, metadataGateway)

	h := grpc_handler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
