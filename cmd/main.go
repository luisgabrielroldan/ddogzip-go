package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"ddogzip/pkg/config"
	"ddogzip/pkg/server"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logOutput := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Logger = logOutput.With().Logger()

	config := config.LoadConfig()

	srv := server.NewServer(config)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Error().Err(err).Msg("Server failed to start")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Received shutdown signal")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("Error during shutdown")
		os.Exit(1)
	}
}
