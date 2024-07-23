package server

import (
	"fmt"
	"log"

	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
	zlog "github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

func NewZipkinReporter(config *config.AppConfig) reporter.Reporter {
	url := fmt.Sprintf("%s://%s:%s/api/v2/spans", config.ZipkinProtocol, config.ZipkinHost, config.ZipkinPort)

	stdLogger := log.New(zlog.Logger, "", log.LstdFlags)
	httpLogger := http.Logger(stdLogger)

	return http.NewReporter(url, httpLogger)
}
