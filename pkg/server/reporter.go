package server

import (
	"fmt"
	"log"

	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
	zlog "github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

// zerologWriter adapts zerolog.Logger to io.Writer
type zerologWriter struct{}

func (w zerologWriter) Write(p []byte) (n int, err error) {
	// Remove trailing newline if present (log.Logger adds it)
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	zlog.Debug().Msg(msg)
	return len(p), nil
}

func NewZipkinReporter(config *config.AppConfig) reporter.Reporter {
	url := fmt.Sprintf("%s://%s:%s/api/v2/spans", config.ZipkinProtocol, config.ZipkinHost, config.ZipkinPort)

	writer := zerologWriter{}
	stdLogger := log.New(writer, "", 0)
	httpLogger := http.Logger(stdLogger)

	return http.NewReporter(url, httpLogger)
}
