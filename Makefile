BINARY_NAME=ddogzip

build:
	go build -o ./bin/$(BINARY_NAME) cmd/main.go

clean:
	go clean
	rm -f ./bin/*

test:
	go test ./...

run:
	go run cmd/main.go

image:
	docker build -t ddogzip .

