BINARY_NAME=ddogzip

build:
	go build -o ./bin/$(BINARY_NAME) cmd/main.go

clean:
	go clean
	rm -f ./bin/*

image:
	docker build -t ddogzip .

