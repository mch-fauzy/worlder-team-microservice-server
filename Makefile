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
	cd microservice-a && swag init -g routes/routes.go
	@echo "Generating swagger docs for microservice-b..."
	cd microservice-b && swag init -g routes/routes.go
	@echo "Swagger documentation generated:"
	@echo "  - microservice-a/docs/"
	@echo "  - microservice-b/docs/"

# Docker operations
run: swagger proto ## Start all services with Docker Compose
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

# Docker rebuild operations
rebuild: swagger proto ## Rebuild and restart all services
	@echo "Rebuilding and restarting all services..."
	docker-compose down
	docker-compose up --build -d
	@echo "All services rebuilt and started."

rebuild-generators: swagger proto ## Rebuild only microservice-a generators
	@echo "Rebuilding microservice-a generators..."
	docker-compose stop microservice-a-generator-temperature microservice-a-generator-humidity microservice-a-generator-pressure microservice-a-generator-light microservice-a-generator-motion
	docker-compose up --build microservice-a-generator-temperature microservice-a-generator-humidity microservice-a-generator-pressure microservice-a-generator-light microservice-a-generator-motion -d
	@echo "Generators rebuilt and started."

rebuild-storage: swagger proto ## Rebuild only microservice-b storage service
	@echo "Rebuilding storage service..."
	docker-compose stop microservice-b
	docker-compose up --build microservice-b -d
	@echo "Storage service rebuilt and started."

# Logging
logs: ## View logs from all services
	docker-compose logs -f

# Quick start for demo
quick-start: run ## Quick start for demo
	@echo "Quick start completed! Services are running."
	@echo ""
	@echo "Storage Control:"
	@echo "  Microservice B (Storage): http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "Individual Generator Control:"
	@echo "  Temperature Generator: http://localhost:8081/swagger/index.html"
	@echo "  Humidity Generator:    http://localhost:8082/swagger/index.html"
	@echo "  Pressure Generator:    http://localhost:8083/swagger/index.html"
	@echo "  Light Generator:       http://localhost:8084/swagger/index.html"
	@echo "  Motion Generator:      http://localhost:8085/swagger/index.html"
	@echo ""
