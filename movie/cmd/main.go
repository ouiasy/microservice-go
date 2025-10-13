package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ouiasy/microservice-go/movie/internal/controller/movie"
	metadatagateway "github.com/ouiasy/microservice-go/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/ouiasy/microservice-go/movie/internal/gateway/rating/http"
	httphandler "github.com/ouiasy/microservice-go/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("http://localhost:8081")
	ratingGateway := ratinggateway.New("http://localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("GET /movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
