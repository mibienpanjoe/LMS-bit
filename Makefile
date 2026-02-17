APP_NAME := lms

.PHONY: run build test test-race bench lint tidy fmt

run:
	go run ./cmd/lms

build:
	go build -o bin/$(APP_NAME) ./cmd/lms

test:
	go test ./...

test-race:
	go test -race ./...

bench:
	go test -run ^$ -bench . -benchmem ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

fmt:
	go fmt ./...
