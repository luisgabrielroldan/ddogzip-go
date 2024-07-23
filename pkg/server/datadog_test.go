package server

import (
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
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
