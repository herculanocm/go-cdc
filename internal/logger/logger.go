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

	switch cfg.AppLogLevel {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Info().Msg("Logger initialized")
	log.Info().Msgf("Log level set to: %s.", logLevel.String())
}
