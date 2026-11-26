.PHONY: build run test tidy fmt

build:
	go build ./...

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .
