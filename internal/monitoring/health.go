// internal/monitoring/health.go
package monitoring

import (
	"context"
	"time"

	"go-cdc/internal/config"
	"go-cdc/internal/runtime"
	"go-cdc/static"

	"github.com/rs/zerolog/log"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context, cfg *config.Config) static.ErrorUtil
}

type HealthMonitor struct {
	cfg           *config.Config
	healthChecker HealthChecker
	metadata      *runtime.Metadata // Injetado via construtor
}

// NewHealthMonitor cria um monitor com todas as dependências injetadas
func NewHealthMonitor(cfg *config.Config, checker HealthChecker, metadata *runtime.Metadata) *HealthMonitor {
	return &HealthMonitor{
		cfg:           cfg,
		healthChecker: checker,
		metadata:      metadata,
	}
}

func (h *HealthMonitor) Start(ctx context.Context) {
	log.Info().
		Dur("interval", time.Duration(h.cfg.HealthCheckIntervalSeconds)*time.Second).
		Msg("Health monitor started")

	ticker := time.NewTicker(time.Duration(h.cfg.HealthCheckIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.check(ctx)
		case <-ctx.Done():
			log.Info().Msg("Health monitor stopped")
			return
		}
	}
}

func (h *HealthMonitor) check(ctx context.Context) {
	start := time.Now()

	errHealthCheck := h.healthChecker.HealthCheck(ctx, h.cfg)
	duration := time.Since(start)

	// Log estruturado com metadados do pod
	logEvent := log.Info()
	if errHealthCheck != nil {
		logEvent = log.Error().Err(errHealthCheck)
	}

	logEvent.
		Str("check_type", "database").
		Str("db_host", h.cfg.DBHost).
		Str("db_name", h.cfg.DBName).
		Float64("duration_seconds", duration.Seconds()).
		Int64("duration_ms", duration.Milliseconds()).
		Bool("success", errHealthCheck == nil).
		Msg("Health check completed")
}
