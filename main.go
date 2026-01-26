// main.go
package main

import (
	"context"
	"go-cdc/database"
	"go-cdc/internal/config"
	"go-cdc/internal/logger"
	"go-cdc/internal/monitoring"
	"go-cdc/internal/runtime"
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

	// 2. Inicializar metadados do runtime (pod/container) - INJEÇÃO
	metadata := runtime.NewMetadata(cfg.AppVersion, cfg.AppEnv)

	// 3. Inicializar logger com contexto global - INJEÇÃO
	dlog.Print("Initializing logger...")
	logger.Init(cfg, metadata.ToLogFields())
	log.Info().Msg("Logger initialized with pod metadata")
	log.Info().Msgf("Configuration: %s", cfg.ToString(false))

	// 4. Inicializar database - INJEÇÃO
	log.Info().Msg("Initializing database connection pool...")
	dbManager, dbErr := database.Init(cfg)
	if dbErr != nil {
		log.Fatal().Err(dbErr).Caller().Msgf("Failed to initialize database: %s", dbErr.ToString())
	}
	defer dbManager.Close()

	// 5. Context para shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().Msg("Initializing application components...")

	// 6. Iniciar monitoramento de saúde - INJEÇÃO DE 3 DEPENDÊNCIAS
	log.Info().Msgf("HealthCheck interval loop: %d seconds", cfg.HealthCheckIntervalSeconds)
	healthMonitor := monitoring.NewHealthMonitor(
		cfg,       // Dependência 1: Config
		dbManager, // Dependência 2: HealthChecker (DatabaseManager implementa a interface)
		metadata,  // Dependência 3: Runtime Metadata
	)
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
