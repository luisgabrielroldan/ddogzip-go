package server

import (
	"errors"
	"log"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
	zlog "github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

var zipkinReporter reporter.Reporter

func zipkinInitReporter(config *config.AppConfig) {
	url := config.ZipkinProtocol + "://" + config.ZipkinHost + ":" + config.ZipkinPort + "/api/v2/spans"

	stdLogger := log.New(zlog.Logger, "", log.LstdFlags)
	httpLogger := http.Logger(stdLogger)

	zipkinReporter = http.NewReporter(url, httpLogger)
}

func zipkinSendSpans(spans []*model.SpanModel) {
	for _, span := range spans {
		zipkinReporter.Send(*span)
	}
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
