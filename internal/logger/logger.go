package logger

import (
	"go-cdc/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.Config) {
	var logLevel zerolog.Level = zerolog.InfoLevel // Nível padrão é Info
	if cfg.AppEnv == "prd" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		logLevel = zerolog.WarnLevel
	}

	if cfg.AppEnv != "prd" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Info().Msg("Logger initialized.")
}
