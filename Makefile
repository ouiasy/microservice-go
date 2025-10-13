.PHONY: run-servers
run-servers:
	docker compose up -d



.PHONY: proto-all-go
proto-all-go: proto-metadata-go proto-movie-go proto-rating-go

.PHONY: proto-metadata-go
proto-metadata-go:
	protoc --go_out=./api/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./api/v1 \
	       api/v1/metadata.proto

.PHONY: proto-movie-go
proto-movie-go:
	protoc --go_out=./api/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./api/v1 \
	       api/v1/movie.proto

.PHONY: proto-rating-go
proto-rating-go:
	protoc --go_out=./api/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./api/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./api/v1 \
	       api/v1/rating.proto

