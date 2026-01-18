package config

import (
	dlog "log"
	"os"
	"strconv"
	"strings"

	viper "github.com/spf13/viper"
)

func ReadEnvFromFileEnabled() bool {
	v, ok := os.LookupEnv("APP_GO_CDC_READ_ENV_FROM_FILE")
	if !ok {
		return false
	}
	v = strings.TrimSpace(v)
	if b, err := strconv.ParseBool(v); err == nil {
		return b
	}
	switch strings.ToLower(v) {
	case "1", "yes", "y", "true", "on":
		return true
	default:
		return false
	}
}

type Config struct {
	AppName string `mapstructure:"APP_GO_CDC_NAME"`
	AppEnv  string `mapstructure:"APP_GO_CDC_ENV"`

	AppReadEnvFromFile string `mapstructure:"APP_GO_CDC_READ_ENV_FROM_FILE"`

	DBHost string `mapstructure:"APP_GO_CDC_DB_HOST"`
	DBPort string `mapstructure:"APP_GO_CDC_DB_PORT"`
	DBUser string `mapstructure:"APP_GO_CDC_DB_USER"`
	DBPass string `mapstructure:"APP_GO_CDC_DB_PASS"`
	DBName string `mapstructure:"APP_GO_CDC_DB_NAME"`
}

func LoadConfig(path string) (*Config, error) {
	dlog.Print("loading config...")

	if ReadEnvFromFileEnabled() {
		dlog.Print("...from file")
		viper.AddConfigPath(path)
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				dlog.Fatalf("config file not found in path: %s", path)
				return nil, err
			}
		}

		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			dlog.Fatalf("unable to decode into struct: %v", err)
			return nil, err
		}
		return &cfg, nil
	}

	dlog.Print("...from environment variables\n")
	var cfg Config
	cfg.AppEnv = os.Getenv("APP_GO_CDC_ENV")
	cfg.AppReadEnvFromFile = os.Getenv("APP_GO_CDC_READ_ENV_FROM_FILE")
	cfg.AppName = os.Getenv("APP_GO_CDC_NAME")
	cfg.DBHost = os.Getenv("APP_GO_CDC_DB_HOST")
	cfg.DBPort = os.Getenv("APP_GO_CDC_DB_PORT")
	cfg.DBUser = os.Getenv("APP_GO_CDC_DB_USER")
	cfg.DBPass = os.Getenv("APP_GO_CDC_DB_PASS")
	cfg.DBName = os.Getenv("APP_GO_CDC_DB_NAME")
	return &cfg, nil
}
