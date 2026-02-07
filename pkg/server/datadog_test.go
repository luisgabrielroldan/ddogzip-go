package server

import (
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

func TestTranslateDDSpanToZipkinSpan(t *testing.T) {
	ddSpan := &DDSpan{
		TraceID:  1,
		SpanID:   1,
		Name:     "test-span",
		Start:    0,
		Duration: uint64(time.Duration(time.Second).Nanoseconds()),
		Service:  "test-service",
		Resource: "test-resource",
		Type:     "test-type",
		Meta: map[string]interface{}{
			"key": "value",
		},
		Error: 0,
	}

	zSpan := translateDDSpanToZipkinSpan(ddSpan)

	assert.Equal(t, model.ID(1), zSpan.ID)
	assert.Equal(t, "test-span", zSpan.Name)
	assert.Equal(t, time.Unix(0, 0), zSpan.Timestamp)
	assert.Equal(t, time.Second, zSpan.Duration)
	assert.Equal(t, "test-service", zSpan.LocalEndpoint.ServiceName)
	assert.Equal(t, "value", zSpan.Tags["key"])
	assert.Equal(t, "test-resource", zSpan.Tags["dd.resource"])
	assert.Equal(t, "test-type", zSpan.Tags["dd.type"])
}

func TestTranslateDDSpanToZipkinSpanWithError(t *testing.T) {
	ddSpan := &DDSpan{
		TraceID:  1,
		SpanID:   1,
		Name:     "test-span",
		Start:    0,
		Duration: uint64(time.Duration(time.Second).Nanoseconds()),
		Service:  "test-service",
		Resource: "test-resource",
		Type:     "test-type",
		Meta: map[string]interface{}{
			"error.msg":   "test error",
			"error.type":  "RuntimeError",
			"error.stack": "stack trace here",
		},
		Error: 1,
	}

	zSpan := translateDDSpanToZipkinSpan(ddSpan)

	assert.NotNil(t, zSpan.Err)
	assert.Equal(t, "test error", zSpan.Err.Error())
	assert.Equal(t, "RuntimeError", zSpan.Tags["error.type"])
	assert.Equal(t, "stack trace here", zSpan.Tags["error.stack"])
}

func TestDecodeDDTraceDataMsgpack(t *testing.T) {
	// Create a valid DDTrace payload
	trace := DDTrace{
		{
			TraceID:  123,
			SpanID:   456,
			Name:     "test-operation",
			Start:    uint64(time.Now().UnixNano()),
			Duration: uint64(time.Millisecond * 100),
			Service:  "test-service",
			Resource: "test-resource",
			Type:     "web",
			Meta:     map[string]interface{}{"key": "value"},
			Metrics:  Metrics{"sample_rate": 1.0},
		},
	}

	data := []DDTrace{trace}
	payload, err := msgpack.Marshal(data)
	assert.NoError(t, err)

	// Decode it
	decoded, err := decodeDDTraceDataMsgpack(payload)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Len(t, *decoded, 1)
	assert.Len(t, (*decoded)[0], 1)

	span := (*decoded)[0][0]
	assert.Equal(t, uint64(123), span.TraceID)
	assert.Equal(t, uint64(456), span.SpanID)
	assert.Equal(t, "test-operation", span.Name)
	assert.Equal(t, "test-service", span.Service)
}

func TestDecodeDDTraceDataUnsupportedVersion(t *testing.T) {
	payload := []byte("test")
	_, err := decodeDDTraceData("v0.1", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol version")
}

func TestDecodeDDTraceDataMalformedPayload(t *testing.T) {
	malformedPayload := []byte("not valid msgpack")
	_, err := decodeDDTraceDataMsgpack(malformedPayload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal payload")
}

func TestDecodeDDTraceDataSupportedVersions(t *testing.T) {
	trace := DDTrace{
		{
			TraceID:  1,
			SpanID:   2,
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
	payload, err := msgpack.Marshal(data)
	assert.NoError(t, err)

	// Test all supported versions
	versions := []string{"v0.3", "v0.4", "v0.5"}
	for _, version := range versions {
		decoded, err := decodeDDTraceData(version, payload)
		assert.NoError(t, err, "Version %s should be supported", version)
		assert.NotNil(t, decoded)
	}
}

func TestDdTraceDataToZipkinSpans(t *testing.T) {
	trace1 := DDTrace{
		{
			TraceID:  1,
			SpanID:   1,
			Name:     "span1",
			Start:    0,
			Duration: 1000,
			Service:  "service1",
			Resource: "resource1",
			Type:     "web",
			Meta:     map[string]interface{}{},
			Metrics:  Metrics{},
		},
		{
			TraceID:  1,
			SpanID:   2,
			ParentID: func() *uint64 { id := uint64(1); return &id }(),
			Name:     "span2",
			Start:    0,
			Duration: 500,
			Service:  "service1",
			Resource: "resource2",
			Type:     "db",
			Meta:     map[string]interface{}{},
			Metrics:  Metrics{},
		},
	}

	trace2 := DDTrace{
		{
			TraceID:  2,
			SpanID:   3,
			Name:     "span3",
			Start:    0,
			Duration: 2000,
			Service:  "service2",
			Resource: "resource3",
			Type:     "web",
			Meta:     map[string]interface{}{},
			Metrics:  Metrics{},
		},
	}

	data := []DDTrace{trace1, trace2}
	spans := ddTraceDataToZipkinSpans(&data)

	assert.Len(t, spans, 3) // 2 spans in trace1 + 1 span in trace2
	assert.Equal(t, "span1", spans[0].Name)
	assert.Equal(t, "span2", spans[1].Name)
	assert.Equal(t, "span3", spans[2].Name)
	assert.NotNil(t, spans[1].ParentID)
	assert.Equal(t, model.ID(1), *spans[1].ParentID)
}
