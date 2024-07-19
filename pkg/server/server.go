package server

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

type Server struct {
	config *config.AppConfig
}

func makeAgentHandler() *http.ServeMux {
	mux := http.NewServeMux()
	log.Info().Msg("Creating new agent handler")

	mux.HandleFunc("/{version}/traces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decoded, err := decodeDDTraceData(body)

		if err != nil {
			log.Error().Err(err).Msg("Failed to decode trace data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		zipkinSpans := ddTraceDataToZipkinSpans(decoded)

		log.Info().Msg("New trace data received")

		zipkinSendSpans(zipkinSpans)

		log.Info().Msgf("Sent %d spans to Zipkin", len(zipkinSpans))

		w.WriteHeader(http.StatusAccepted)
	})

	return mux
}

func (s *Server) Start() {
	config := s.config

	log.Info().Msgf("Server listening on %s", config.ListenAddr)

	zipkinInitReporter(config)

	err := http.ListenAndServe(s.config.ListenAddr, makeAgentHandler())

	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("Server closed")
	} else if err != nil {
		log.Error().Err(err).Msg("An error occurred")
		os.Exit(1)
	}
}

func NewServer(config *config.AppConfig) *Server {
	return &Server{
		config: config,
	}
}
