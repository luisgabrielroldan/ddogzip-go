package main

import (
	"os"
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

	server := server.NewServer(config)

	server.Start()
}
