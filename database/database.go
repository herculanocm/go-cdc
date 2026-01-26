// database/database.go
package database

import (
	"context"
	"go-cdc/database/sqlserver"
	"go-cdc/internal/config"
	"go-cdc/static"

	"github.com/rs/zerolog/log"
)

// DatabaseManager gerencia a conexão com o banco
type DatabaseManager struct {
	sqlServer *sqlserver.SQLServer
}

func Init(config *config.Config) (*DatabaseManager, static.ErrorUtil) {
	log.Info().Msg("Initializing database configuration...")

	log.Info().Msg("Checking configuration environments...")
	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPass == "" || config.DBName == "" {
		log.Error().Caller().Msg("Database configuration is incomplete.")
		return nil, static.ErrEnvVarMissing
	}

	log.Info().Msgf("Database technology: %s", config.DBTecnology)
	switch config.DBTecnology {
	case "sqlserver":
		log.Info().Msg("Initializing SQL Server package...")
		sqlServer, errDb := sqlserver.NewSQLServer(config)
		if errDb != nil {
			return nil, errDb
		}

		log.Info().Msg("Database initialized successfully.")
		return &DatabaseManager{sqlServer: sqlServer}, nil

	default:
		log.Error().Caller().Msgf("Unsupported database technology: %s", config.DBTecnology)
		return nil, static.ErrUnsupportedDBTechnology
	}
}

// HealthCheck implementa a interface HealthChecker
func (dm *DatabaseManager) HealthCheck(ctx context.Context, cfg *config.Config) static.ErrorUtil {
	log.Info().Msg("Performing database health check...")

	errDb := dm.sqlServer.HealthCheck(ctx, cfg.DBPingTimeoutSeconds)
	if errDb != nil {
		log.Error().Caller().Err(errDb).Msg("Database health check failed")
		return errDb
	}

	log.Info().Msg("Database health check succeeded.")
	return nil
}

// GetSQLServer retorna a instância do SQL Server (se precisar acessar diretamente)
func (dm *DatabaseManager) GetSQLServer() *sqlserver.SQLServer {
	return dm.sqlServer
}

// Close fecha todas as conexões
func (dm *DatabaseManager) Close() error {
	if dm.sqlServer != nil {
		return dm.sqlServer.Close()
	}
	return nil
}
