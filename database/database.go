package database

import (
	"go-cdc/internal/config"

	"github.com/rs/zerolog/log"
)

type connParams struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func (c *connParams) GetConnString(showPassword bool) string {
	password := "******"
	if showPassword {
		password = c.DBPass
	}

	return "server=" + c.DBHost + ";" +
		"user id=" + c.DBUser + ";" +
		"password=" + password + ";" +
		"port=" + c.DBPort + ";" +
		"database=" + c.DBName + ";" +
		"encrypt=false"
}

var ConnConfig *connParams

func Init(config *config.Config) error {
	log.Info().Msg("Initializing database configuration...")

	ConnConfig = &connParams{
		DBHost: config.DBHost,
		DBPort: config.DBPort,
		DBUser: config.DBUser,
		DBPass: config.DBPass,
		DBName: config.DBName,
	}

	log.Info().Msgf("Database technology: %s", config.DBTecnology)
	log.Info().Msgf("Database connection string: %s", ConnConfig.GetConnString(false))
	log.Info().Msg("Database configuration initialized.")
	return nil
}
