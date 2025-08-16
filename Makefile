.PHONY: all tidy fmt run test

all: tidy fmt test

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

run:
	go run -race ./cmd/app/main.go --config ./config/config.local.yaml

test:
	go test ./... -v

start-dep:
	mkdir -p ./postgres_data
	docker compose up -d postgres

stop-dep:
	docker compose down postgres

docker-build:
	docker build -f ./docker/Dockerfile -t kart-challenge .

docker-start:
	docker compose up -d postgres
	docker compose up -d --build kart-challenge

docker-stop:
	docker compose down kart-challenge
	docker compose down postgres

help:
	@echo "Available commands:"
	@echo "  make all          - Run tidy, fmt, and test"
	@echo "  make tidy         - Tidy go modules"
	@echo "  make fmt          - Format Go code"
	@echo "  make run          - Run the application with local config"
	@echo "  make test         - Run tests"
	@echo "  make start-dep    - Start PostgreSQL dependency"
	@echo "  make stop-dep     - Stop PostgreSQL dependency"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-start - Start Docker containers"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make help         - Show this help message"