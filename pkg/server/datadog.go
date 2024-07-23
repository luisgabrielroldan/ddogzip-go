package server

import (
	"errors"
	"time"
	"github.com/vmihailenco/msgpack"
	"github.com/openzipkin/zipkin-go/model"
)

type Metrics map[string]float64

type DDSpan struct {
	SpanID   uint64                 `msgpack:"span_id"`
	TraceID  uint64                 `msgpack:"trace_id"`
	ParentID *uint64                `msgpack:"parent_id, omitempty"`
	Name     string                 `msgpack:"name"`
	Start    uint64                 `msgpack:"start"`
	Duration uint64                 `msgpack:"duration"`
	Service  string                 `msgpack:"service"`
	Resource string                 `msgpack:"resource"`
	Error    int32                  `msgpack:"error, omitempty"`
	Type     string                 `msgpack:"type"`
	Meta     map[string]interface{} `msgpack:"meta"`
	Metrics  Metrics                `msgpack:"metrics"`
}

type DDTrace []DDSpan

func decodeDDTraceData(payload []byte) (*[]DDTrace, error) {
	var data []DDTrace
	err := msgpack.Unmarshal(payload, &data)
	return &data, err
}

func ddTraceDataToZipkinSpans(data *[]DDTrace) []*model.SpanModel {
	var spans []*model.SpanModel = make([]*model.SpanModel, 0)

	for _, trace := range *data {
		for _, span := range trace {
			spans = append(spans, translateDDSpanToZipkinSpan(&span))
		}
	}

	return spans
}

func translateDDSpanToZipkinSpan(span *DDSpan) *model.SpanModel {
	var zSpan model.SpanModel
	var tags = make(map[string]string)

	zSpan.TraceID.High = 0
	zSpan.TraceID.Low = span.TraceID

	zSpan.ID = model.ID(span.SpanID)

	if span.ParentID != nil {
		parentID := model.ID(*span.ParentID)
		zSpan.ParentID = &parentID
	}

	zSpan.Name = span.Name
	zSpan.Timestamp = time.Unix(0, int64(span.Start))
	zSpan.Duration = time.Duration(span.Duration)

	for key, value := range span.Meta {
		tags[key] = value.(string)
	}

	tags["dd.resource"] = span.Resource
	tags["dd.type"] = span.Type

	zSpan.Tags = tags

	errorMessage := span.Meta["error.msg"]

	if span.Error != 0 && errorMessage != nil {
		zSpan.Err = errors.New(span.Meta["error.msg"].(string))
	}

	zSpan.LocalEndpoint = &model.Endpoint{ServiceName: span.Service}

	return &zSpan
}
