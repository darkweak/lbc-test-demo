.PHONY: build lint prepare run

build:
	go build ./cmd

lint:
	golangci-lint run -v ./...

prepare:
	go install golang.org/dl/go1.26.3@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
	go1.26.3 mod tidy

run:
	go run ./cmd

tests:
	go test -v ./...