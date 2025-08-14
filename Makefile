.PHONY: all tidy fmt run test

all: tidy fmt test

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

run:
	go run -race ./cmd/server/main.go

test:
	go test ./... -race -v