// database/sqlserver/main.go
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

// SQLServer encapsula a conexão e configuração
type SQLServer struct {
	db         *sql.DB
	connConfig *connParams
}

// NewSQLServer cria e inicializa uma nova instância (construtor)
func NewSQLServer(config *config.Config) (*SQLServer, static.ErrorUtil) {
	connConfig := &connParams{
		DBHost: config.DBHost,
		DBPort: config.DBPort,
		DBUser: config.DBUser,
		DBPass: config.DBPass,
		DBName: config.DBName,
	}

	log.Info().Msg("Configuring database connection...")

	db, err := sql.Open("sqlserver", connConfig.GetConnString(true))
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to open database connection")
		return nil, static.NewErrorUtil("Failed to open database connection", "SQLSERVER_INIT_FAILED", err, err.Error())
	}

	// Pool sizing para CDC workload
	db.SetMaxOpenConns(config.DBMaxOpenConns)
	db.SetMaxIdleConns(config.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.DBConnMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(config.DBConnMaxIdleTime) * time.Minute)

	sqlServer := &SQLServer{
		db:         db,
		connConfig: connConfig,
	}

	// Verifica conexão inicial
	log.Info().Msg("Pinging database to verify connection...")
	if err := sqlServer.HealthCheck(context.Background(), config.DBPingTimeoutSeconds); err != nil {
		return nil, err
	}

	return sqlServer, nil
}

// GetDB retorna a instância do pool (thread-safe por design do sql.DB)
func (s *SQLServer) GetDB() *sql.DB {
	return s.db
}

// GetConnConfig retorna a configuração de conexão
func (s *SQLServer) GetConnConfig() *connParams {
	return s.connConfig
}

// HealthCheck verifica a saúde da conexão
func (s *SQLServer) HealthCheck(ctx context.Context, timeoutSeconds int) static.ErrorUtil {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		log.Error().Caller().Err(err).Msg("Database health check failed")
		return static.NewErrorUtil("Database health check failed", "SQLSERVER_HEALTH_CHECK_FAILED", err, err.Error())
	}

	return nil
}

// Close fecha a conexão com o banco
func (s *SQLServer) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
