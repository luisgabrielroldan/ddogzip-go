package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
)

type Server struct {
	config         *config.AppConfig
	zipkinReporter reporter.Reporter
	httpServer     *http.Server
}

func NewServer(config *config.AppConfig) *Server {
	return &Server{
		config:         config,
		zipkinReporter: NewZipkinReporter(config),
	}
}

func (s *Server) Start() error {
	config := s.config

	s.httpServer = &http.Server{
		Addr:    config.ListenAddr,
		Handler: makeAgentHandler(s),
	}

	log.Info().Msgf("Server listening on %s", config.ListenAddr)

	err := s.httpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("Server closed")
		return nil
	} else if err != nil {
		log.Error().Err(err).Msg("An error occurred")
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var shutdownErr error

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error shutting down HTTP server")
		shutdownErr = err
	}

	// Always close reporter to flush pending spans
	if err := s.zipkinReporter.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing Zipkin reporter")
		if shutdownErr == nil {
			shutdownErr = err
		}
	}

	log.Info().Msg("Server shutdown complete")
	return shutdownErr
}

func (s *Server) reportSpans(spans []*model.SpanModel) {
	for _, span := range spans {
		s.zipkinReporter.Send(*span)
	}
}

func makeAgentHandler(server *Server) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/{version}/traces", func(w http.ResponseWriter, r *http.Request) {
		// Support both PUT and POST methods
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
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
