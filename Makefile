.PHONY: run-servers
run-servers:
	docker compose up -d



.PHONY: proto-all-go
proto-all-go: proto-metadata-go proto-movie-go proto-rating-go

.PHONY: proto-metadata-go
proto-metadata-go:
	protoc --go_out=./_proto/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./_proto/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./_proto/v1 \
	       _proto/v1/metadata.proto

.PHONY: proto-movie-go
proto-movie-go:
	protoc --go_out=./_proto/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./_proto/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./_proto/v1 \
	       _proto/v1/movie.proto

.PHONY: proto-rating-go
proto-rating-go:
	protoc --go_out=./_proto/gen/go --go_opt=paths=source_relative \
	       --go-grpc_out=./_proto/gen/go --go-grpc_opt=paths=source_relative \
	       -I=./_proto/v1 \
	       _proto/v1/rating.proto

.PHONY: install_deps
install_deps:
	go get -u google.golang.org