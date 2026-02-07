FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o ddogzip cmd/main.go


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /app/ddogzip .

ENTRYPOINT ["./ddogzip"]
