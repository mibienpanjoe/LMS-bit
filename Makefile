APP_NAME := lms

.PHONY: run build test test-race lint tidy fmt

run:
	go run ./cmd/lms

build:
	go build -o bin/$(APP_NAME) ./cmd/lms

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

fmt:
	go fmt ./...
