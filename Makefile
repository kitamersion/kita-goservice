.PHONY: build run test clean docker-up docker-down kafka-format

# Build all applications
build:
	go build -o bin/api ./cmd/api
	go build -o bin/consumer ./cmd/consumer

# Run API server
run-api:
	go run ./cmd/api

# Run consumer
run-consumer:
	go run ./cmd/consumer

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

# Development
dev-api:
	air -c .air.toml ./cmd/api

dev-consumer:
	air -c .air.toml ./cmd/consumer

