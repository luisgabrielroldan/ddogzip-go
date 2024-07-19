package server

import (
	"github.com/vmihailenco/msgpack"
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
