package main

import (
	"context"
	"go-cdc/database"
	"go-cdc/internal/config"
	"go-cdc/internal/logger"
	"go-cdc/internal/monitoring"
	dlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	dlog.Print("Starting go-cdc application...")

	// 1. Carregar configuração
	cfg, err := config.LoadConfig(".")
	if err != nil {
		dlog.Fatalf("Failed to load config: %v", err.ToString())
	}

	// 2. Inicializar logger
	dlog.Print("Initializing logger...")
	logger.Init(cfg)
	log.Info().Msg("Logger initialized")
	log.Info().Msgf("Configuration: %s", cfg.ToString(false))

	// 3. Inicializar database
	log.Info().Msg("Initializing database connection pool...")
	dbErr := database.Init(cfg)
	if dbErr != nil {
		log.Fatal().Err(dbErr).Caller().Msgf("Failed to initialize database: %s", dbErr.ToString())
	}

	// 4. Context para shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().Msg("Initializing application components...")
	// 5. Iniciar monitoramento de saúde
	log.Info().Msgf("HealthCheck interval loop: %d seconds", cfg.HealthCheckIntervalSeconds)
	healthMonitor := monitoring.NewHealthMonitor(time.Duration(cfg.HealthCheckIntervalSeconds) * time.Second)
	go healthMonitor.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Shutdown signal received, exiting...")
	cancel()

	// Aguarda um pouco para garantir log final
	time.Sleep(500 * time.Millisecond)

	log.Info().Msg("Application stopped")
}
