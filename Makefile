.PHONY: generate build run test clean docker-up docker-down proto

# Generate GraphQL code
generate:
	go mod tidy
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate

# Build all applications
build: generate
	go build -o bin/api ./cmd/api
	go build -o bin/consumer ./cmd/consumer
	go build -o bin/graph ./cmd/graph

# Run API server
run-api: generate
	go run ./cmd/api

# Run consumer
run-consumer: generate
	go run ./cmd/consumer

# Run the application
run: generate
	go run cmd/api/main.go

# Run the GraphQL server
run-graph: generate
	go run cmd/graph/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-build:
	docker-compose build


proto:
	protoc -I=./proto --go_out=./internal/events/ ./proto/**/**/*.proto

# Initialize GraphQL (run this once)
init-graphql:
	go run github.com/99designs/gqlgen init
