package database

import (
	"context"
	"database/sql"
	"go-cdc/internal/config"
	"go-cdc/static"
	"time"

	"github.com/rs/zerolog/log"

	// Importa o driver oficial
	_ "github.com/microsoft/go-mssqldb"
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
var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func Connect(config *config.Config) error {
	var err error
	log.Info().Msg("Connecting to the database...")
	db, err = sql.Open("sqlserver", ConnConfig.GetConnString(true))
	if err != nil {
		log.Error().Err(err).Msg("Error opening database connection")
		return err
	}

	// Testa a conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error pinging database")
		return err
	}

	log.Info().Msg("Database connection established.")
	return nil
}

func Init(config *config.Config) static.ErrorUtil {
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

	log.Info().Msg("Checking configuration environments...")
	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPass == "" || config.DBName == "" {
		log.Error().Msg("Database configuration is incomplete.")
		return static.ErrEnvVarMissing
	}

	return nil
}
