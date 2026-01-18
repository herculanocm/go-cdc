package sqlserver

import (
	"context"
	"database/sql"
	"go-cdc/internal/config"
	"go-cdc/static"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/rs/zerolog/log"
)

var db *sql.DB
var ConnConfig *connParams

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
		"encrypt=" + "true" + ";" +
		"trustservercertificate=" + "true" + ";"
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
	ConnConfig = &connParams{
		DBHost: config.DBHost,
		DBPort: config.DBPort,
		DBUser: config.DBUser,
		DBPass: config.DBPass,
		DBName: config.DBName,
	}

	err := Connect(config)
	if err != nil {
		return static.NewErrorUtil("Failed to initialize SQL Server database", "SQLSERVER_INIT_FAILED", err, err.Error())
	}
	return nil
}
