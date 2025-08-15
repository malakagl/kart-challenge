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