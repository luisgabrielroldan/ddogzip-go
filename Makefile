BINARY_NAME=ddogzip
GOLANGCI_LINT_VERSION=v1.64.8
DOCKER_IMAGE=luisgabrielroldan/ddogzip

.PHONY: build clean test lint run image image-push check

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
	docker build -t $(DOCKER_IMAGE) .

image-multiarch:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE) .

image-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):latest --push .

check: test lint build

