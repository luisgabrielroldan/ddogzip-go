package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"

	"ddogzip/pkg/config"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestHealthEndpointMethodNotAllowed(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestTracesEndpointMethodNotAllowed(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	req := httptest.NewRequest(http.MethodGet, "/v0.4/traces", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestTracesEndpointInvalidPayload(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	req := httptest.NewRequest(http.MethodPut, "/v0.4/traces", bytes.NewReader([]byte("invalid msgpack")))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTracesEndpointUnsupportedVersion(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	trace := DDTrace{
		{
			TraceID:  1,
			SpanID:   1,
			Name:     "test",
			Start:    0,
			Duration: 1000,
			Service:  "test",
			Resource: "test",
			Type:     "web",
			Meta:     map[string]interface{}{},
			Metrics:  Metrics{},
		},
	}

	data := []DDTrace{trace}
	payload, _ := msgpack.Marshal(data)

	req := httptest.NewRequest(http.MethodPut, "/v0.1/traces", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTracesEndpointSuccessPUT(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	trace := DDTrace{
		{
			TraceID:  1,
			SpanID:   1,
			Name:     "test-operation",
			Start:    0,
			Duration: 1000,
			Service:  "test-service",
			Resource: "test-resource",
			Type:     "web",
			Meta:     map[string]interface{}{"key": "value"},
			Metrics:  Metrics{},
		},
	}

	data := []DDTrace{trace}
	payload, _ := msgpack.Marshal(data)

	req := httptest.NewRequest(http.MethodPut, "/v0.4/traces", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestTracesEndpointSuccessPOST(t *testing.T) {
	cfg := &config.AppConfig{
		ListenAddr:     ":8126",
		ZipkinProtocol: "http",
		ZipkinHost:     "localhost",
		ZipkinPort:     "9411",
	}

	srv := NewServer(cfg)
	handler := makeAgentHandler(srv)

	trace := DDTrace{
		{
			TraceID:  1,
			SpanID:   1,
			Name:     "test-operation",
			Start:    0,
			Duration: 1000,
			Service:  "test-service",
			Resource: "test-resource",
			Type:     "web",
			Meta:     map[string]interface{}{"key": "value"},
			Metrics:  Metrics{},
		},
	}

	data := []DDTrace{trace}
	payload, _ := msgpack.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/v0.5/traces", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}
