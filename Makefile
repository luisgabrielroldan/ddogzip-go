BINARY_NAME=ddogzip
GOLANGCI_LINT_VERSION=v1.64.8

.PHONY: build clean test lint run image check

build:
	go build -o ./bin/$(BINARY_NAME) cmd/main.go

clean:
	go clean
	rm -f ./bin/*

test:
	go test ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

run:
	go run cmd/main.go

image:
	docker build -t ddogzip .

check: test lint build

