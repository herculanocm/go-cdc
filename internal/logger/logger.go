// internal/logger/logger.go
package logger

import (
	"go-cdc/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.Config) {
	var logLevel zerolog.Level = zerolog.InfoLevel

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

	// Adiciona campos globais ao contexto (pod metadata)
	if cfg.ToLogFields() != nil && len(cfg.ToLogFields()) > 0 {
		ctx := log.Logger.With()
		for key, value := range cfg.ToLogFields() {
			ctx = ctx.Interface(key, value)
		}
		log.Logger = ctx.Logger()
	}

	log.Info().Msg("Logger initialized")
	log.Info().Msgf("Log level set to: %s", logLevel.String())
}
