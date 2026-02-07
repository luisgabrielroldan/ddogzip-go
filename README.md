# DDogZip

[![Docker Image](https://img.shields.io/docker/v/luisgabrielroldan/ddogzip?label=Docker%20Hub)](https://hub.docker.com/r/luisgabrielroldan/ddogzip)
[![Docker Pulls](https://img.shields.io/docker/pulls/luisgabrielroldan/ddogzip)](https://hub.docker.com/r/luisgabrielroldan/ddogzip)

DDogZip is a proxy server that receives traces from Datadog-instrumented applications and forwards them to a Zipkin collector. This allows you to use Datadog's APM instrumentation while sending trace data to Zipkin for visualization and analysis.

This tool is ideal for local development environments where you want to debug tracing without sending data to Datadog's infrastructure.

## Features

- Receives traces via Datadog Agent protocol (v0.3, v0.4, v0.5)
- Converts Datadog traces to Zipkin format
- Supports both PUT and POST HTTP methods
- Graceful shutdown with span flushing
- Health check endpoint for container orchestration
- Enhanced error capture (error.msg, error.type, error.stack)
- Multi-architecture Docker images (amd64, arm64)

## Configuration

DDogZip can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8126` | Address and port to listen on |
| `ZIPKIN_PROTOCOL` | `http` | Protocol for Zipkin endpoint (http/https) |
| `ZIPKIN_HOST` | `localhost` | Zipkin collector hostname |
| `ZIPKIN_PORT` | `9411` | Zipkin collector port |

## Usage

### Docker (Recommended)

```bash
docker run -p 8126:8126 \
  -e ZIPKIN_HOST=zipkin \
  -e ZIPKIN_PORT=9411 \
  luisgabrielroldan/ddogzip
```

### Running Locally

```bash
# Using default configuration
go run cmd/main.go

# With custom configuration
LISTEN_ADDR=:9000 ZIPKIN_HOST=zipkin.example.com go run cmd/main.go
```

### Building

```bash
make build
./bin/ddogzip
```

### Building Docker Image Locally

```bash
# Build for current platform
make image

# Build multi-arch image
make image-multiarch
```

## Docker Compose Example

Here's a complete example with DDogZip and Zipkin:

```yaml
version: '3.8'

services:
  zipkin:
    image: openzipkin/zipkin:latest
    ports:
      - "9411:9411"
    environment:
      - STORAGE_TYPE=mem

  ddogzip:
    image: luisgabrielroldan/ddogzip:latest
    ports:
      - "8126:8126"
    environment:
      - LISTEN_ADDR=:8126
      - ZIPKIN_PROTOCOL=http
      - ZIPKIN_HOST=zipkin
      - ZIPKIN_PORT=9411
    depends_on:
      - zipkin
```

## API Endpoints

### `PUT /{version}/traces` or `POST /{version}/traces`

Receives Datadog trace data in msgpack format.

- **Supported versions**: `v0.3`, `v0.4`, `v0.5`
- **Content-Type**: `application/msgpack`
- **Response**: `202 Accepted` on success

### `GET /health`

Health check endpoint for container orchestration.

- **Response**: `200 OK` with body "OK"

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean
```

## How It Works

1. Datadog-instrumented applications send traces to DDogZip on port 8126
2. DDogZip decodes the msgpack-encoded Datadog trace data
3. Traces are converted from Datadog format to Zipkin format
4. Converted spans are forwarded to the Zipkin collector

## License

MIT License

