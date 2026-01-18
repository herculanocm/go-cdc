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
		dlog.Fatalf("Failed to load config: %v", err.ToString())
	}

	dlog.Print("Config loaded.")
	dlog.Print("Initializing logger...")
	logger.Init(cfg)

	log.Info().Msgf("Environments: %s", cfg.ToString(false))

	dbErr := database.Init(cfg)
	if dbErr != nil {
		log.Fatal().Err(dbErr).Caller().Msgf("Failed to initialize database: %s", dbErr.ToString())
	}

	log.Info().Msg("Application started.")
}
