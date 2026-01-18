package database

import (
	"go-cdc/database/sqlserver"
	"go-cdc/internal/config"
	"go-cdc/static"

	"github.com/rs/zerolog/log"
)

func Init(config *config.Config) static.ErrorUtil {
	log.Info().Msg("Initializing database configuration...")

	log.Info().Msg("Checking configuration environments...")
	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPass == "" || config.DBName == "" {
		log.Error().Msg("Database configuration is incomplete.")
		return static.ErrEnvVarMissing
	}
	log.Info().Msgf("Database technology: %s", config.DBTecnology)
	switch config.DBTecnology {
	case "sqlserver":
		{
			log.Info().Msg("Testing SQL Server connection...")
			errDb := sqlserver.Init(config)
			if errDb != nil {
				return errDb
			}
		}
	default:
		{
			log.Error().Msgf("Unsupported database technology: %s", config.DBTecnology)
			return static.ErrUnsupportedDBTechnology
		}
	}

	log.Info().Msg("Database initialized successfully.")

	return nil
}
