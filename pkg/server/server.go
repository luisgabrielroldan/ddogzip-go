package server

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

type Server struct {
	config         *config.AppConfig
	zipkinReporter reporter.Reporter
}

func NewServer(config *config.AppConfig) *Server {
	return &Server{
		config:         config,
		zipkinReporter: NewZipkinReporter(config),
	}
}

func (s *Server) Start() {
	config := s.config

	log.Info().Msgf("Server listening on %s", config.ListenAddr)

	err := http.ListenAndServe(s.config.ListenAddr, makeAgentHandler(s))

	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("Server closed")
	} else if err != nil {
		log.Error().Err(err).Msg("An error occurred")
		os.Exit(1)
	}
}

func (s *Server) reportSpans(spans []*model.SpanModel) {
	for _, span := range spans {
		s.zipkinReporter.Send(*span)
	}
}

func makeAgentHandler(server *Server) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/{version}/traces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		version := r.PathValue("version")
		decoded, err := decodeDDTraceData(version, body)

		if err != nil {
			log.Error().Err(err).Msg("Failed to decode trace data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		zipkinSpans := ddTraceDataToZipkinSpans(decoded)

		log.Info().Msg("New trace data received")

		server.reportSpans(zipkinSpans)

		log.Info().Msgf("Sent %d spans to Zipkin", len(zipkinSpans))

		w.WriteHeader(http.StatusAccepted)
	})

	return mux
}
