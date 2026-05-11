.PHONY: build run test docker-build migrate proto

build:
	go build -o bin/vigilante cmd/vigilante/*.go

run: build
	./bin/vigilante serve

test:
	go test ./... -v

docker-build:
	docker build -t vigilante:latest .

migrate: build
	./bin/vigilante migrate

proto:
	protoc --go_out=. --go-grpc_out=. internal/grpc/metrics.proto
