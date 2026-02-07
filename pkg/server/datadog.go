package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/vmihailenco/msgpack"
)

type Metrics map[string]float64

type DDSpan struct {
	SpanID   uint64                 `msgpack:"span_id"`
	TraceID  uint64                 `msgpack:"trace_id"`
	ParentID *uint64                `msgpack:"parent_id,omitempty"`
	Name     string                 `msgpack:"name"`
	Start    uint64                 `msgpack:"start"`
	Duration uint64                 `msgpack:"duration"`
	Service  string                 `msgpack:"service"`
	Resource string                 `msgpack:"resource"`
	Error    int32                  `msgpack:"error,omitempty"`
	Type     string                 `msgpack:"type"`
	Meta     map[string]interface{} `msgpack:"meta"`
	Metrics  Metrics                `msgpack:"metrics"`
}

type DDTrace []DDSpan

func decodeDDTraceData(version string, payload []byte) (*[]DDTrace, error) {
	switch version {
	case "v0.3", "v0.4", "v0.5":
		// v0.3, v0.4, and v0.5 all use the same msgpack encoding for trace data
		return decodeDDTraceDataMsgpack(payload)
	default:
		return nil, fmt.Errorf("unsupported protocol version: %s", version)
	}
}

func decodeDDTraceDataMsgpack(payload []byte) (*[]DDTrace, error) {
	var data []DDTrace
	if err := msgpack.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &data, nil
}

func ddTraceDataToZipkinSpans(data *[]DDTrace) []*model.SpanModel {
	spans := make([]*model.SpanModel, 0, len(*data))

	for _, trace := range *data {
		for _, span := range trace {
			spans = append(spans, translateDDSpanToZipkinSpan(&span))
		}
	}

	return spans
}

func translateDDSpanToZipkinSpan(span *DDSpan) *model.SpanModel {
	zSpan := model.SpanModel{
		SpanContext: model.SpanContext{
			TraceID: model.TraceID{Low: span.TraceID},
			ID:      model.ID(span.SpanID),
		},
		Name:          span.Name,
		Timestamp:     time.Unix(0, int64(span.Start)),
		Duration:      time.Duration(span.Duration),
		LocalEndpoint: &model.Endpoint{ServiceName: span.Service},
		Tags:          make(map[string]string),
	}

	if span.ParentID != nil {
		parentID := model.ID(*span.ParentID)
		zSpan.ParentID = &parentID
	}

	// Add tags from DDSpan.Meta
	for key, value := range span.Meta {
		if strValue, ok := value.(string); ok {
			zSpan.Tags[key] = strValue
		}
	}

	zSpan.Tags["dd.resource"] = span.Resource
	zSpan.Tags["dd.type"] = span.Type

	if span.Error != 0 {
		if errorMessage, ok := span.Meta["error.msg"].(string); ok {
			zSpan.Err = errors.New(errorMessage)
		}

		// Capture error.type if present
		if errorType, ok := span.Meta["error.type"].(string); ok {
			zSpan.Tags["error.type"] = errorType
		}

		// Capture error.stack if present (truncate if too long)
		if errorStack, ok := span.Meta["error.stack"].(string); ok {
			const maxStackLength = 10000 // Limit stack trace length
			if len(errorStack) > maxStackLength {
				errorStack = errorStack[:maxStackLength] + "... (truncated)"
			}
			zSpan.Tags["error.stack"] = errorStack
		}
	}

	return &zSpan
}
