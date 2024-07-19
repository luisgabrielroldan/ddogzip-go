FROM golang:1.22-alpine as builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ddogzip cmd/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates mailcap && addgroup -S app && adduser -S app -G app
USER app
WORKDIR /app
COPY --from=builder /app/ddogzip .
ENTRYPOINT ["./ddogzip"]
