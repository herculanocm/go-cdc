package main

import (
	"go-cdc/internal/config"
	"go-cdc/internal/logger"
	dlog "log"
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
}
