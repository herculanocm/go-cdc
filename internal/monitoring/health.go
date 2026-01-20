package monitoring

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type HealthMonitor struct {
	interval time.Duration
}

func NewHealthMonitor(interval time.Duration) *HealthMonitor {
	return &HealthMonitor{
		interval: interval,
	}
}

func (h *HealthMonitor) Start(ctx context.Context) {
	log.Info().
		Dur("interval", h.interval).
		Msg("Health monitor started")

	ticker := time.NewTicker(h.interval)
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

	log.Info().
		Dur("duration", time.Since(start)).
		Msg("Health check completed")
}
