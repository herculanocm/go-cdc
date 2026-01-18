package main

import (
	"go-cdc/database"
	"go-cdc/internal/config"
	"go-cdc/internal/logger"
	dlog "log"

	"github.com/rs/zerolog/log"
)

func main() {
	// This is a placeholder for the main function.
	dlog.Print("Starting...")

	cfg, err := config.LoadConfig(".")
	if err != nil {
		dlog.Fatalf("Failed to load config: %v", err)
	}

	dlog.Print("Config loaded.")
	dlog.Print("Initializing logger...")
	logger.Init(cfg)

	err = database.Init(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	log.Info().Msg("Application started.")
}
