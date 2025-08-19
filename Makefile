.PHONY: help build run stop clean logs proto swagger

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Protocol Buffers
proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative shared/proto/sensor/sensor.proto

# Swagger documentation
swagger: ## Generate swagger documentation for all microservices
	@echo "Installing swag if not present..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating swagger docs for microservice-a..."
	cd microservice-a && swag init -g cmd/server/main.go
	@echo "Generating swagger docs for microservice-b..."
	cd microservice-b && swag init -g cmd/server/main.go
	@echo "Swagger documentation generated:"
	@echo "  - microservice-a/docs/"
	@echo "  - microservice-b/docs/"

# Build
build: proto ## Build all services
	@echo "Building microservice-a..."
	cd microservice-a && go build -o bin/microservice-a cmd/server/main.go
	@echo "Building microservice-b..."
	cd microservice-b && go build -o bin/microservice-b cmd/server/main.go

# Docker operations
run: ## Start all services with Docker Compose
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started. Use 'make logs' to view logs."

stop: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down

clean: ## Clean up containers and volumes
	@echo "Cleaning up..."
	docker-compose down --volumes --remove-orphans

restart: stop run ## Restart all services

# Logging
logs: ## View logs from all services
	docker-compose logs -f

# Quick start for demo
quick-start: run ## Quick start for demo
	@echo "Quick start completed! Services are running."
	@echo "Visit http://localhost:8081/swagger/index.html for API documentation"
	@echo "Temperature sensor: http://localhost:8080"
	@echo "Humidity sensor: http://localhost:8082" 
	@echo "Pressure sensor: http://localhost:8083"
