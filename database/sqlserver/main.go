package sqlserver

import (
	"context"
	"database/sql"
	"go-cdc/internal/config"
	"go-cdc/static"
	"sync"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/rs/zerolog/log"
)

var once sync.Once
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

// GetDB retorna a instância única do pool (thread-safe)
func GetDB() *sql.DB {
	return db
}

func HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return GetDB().PingContext(ctx)
}

func Init(config *config.Config) static.ErrorUtil {
	ConnConfig = &connParams{
		DBHost: config.DBHost,
		DBPort: config.DBPort,
		DBUser: config.DBUser,
		DBPass: config.DBPass,
		DBName: config.DBName,
	}

	log.Info().Msg("Configuring database connection...")
	var initErr error
	once.Do(func() {
		db, initErr = sql.Open("sqlserver", ConnConfig.GetConnString(true))
		if initErr != nil {
			return
		}

		// Pool sizing para CDC workload
		db.SetMaxOpenConns(config.DBMaxOpenConns) // 25-50 para CDC
		db.SetMaxIdleConns(config.DBMaxIdleConns) // ~10-20
		db.SetConnMaxLifetime(time.Duration(config.DBConnMaxLifetime) * time.Minute)
		db.SetConnMaxIdleTime(5 * time.Minute) // Libera conexões idle

		initErr = db.PingContext(context.Background())
	})

	if initErr != nil {
		initErrUtil := static.NewErrorUtil("Failed to open database connection", "SQLSERVER_INIT_FAILED", initErr, initErr.Error())
		return initErrUtil
	}

	log.Info().Msg("Pinging database to verify connection...")
	err := HealthCheck(context.Background())
	if err != nil {
		return static.NewErrorUtil("Failed to ping database", "SQLSERVER_PING_FAILED", err, err.Error())
	}

	return nil
}
