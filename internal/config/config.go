package config

import (
	"fmt"
	"go-cdc/static"
	dlog "log"
	"net"
	"os"
	"reflect"
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
	AppName    string `mapstructure:"APP_GO_CDC_NAME"`
	AppEnv     string `mapstructure:"APP_GO_CDC_ENV"`
	AppVersion string `mapstructure:"APP_GO_CDC_VERSION"` // Adicione esta linha

	PodHostname string `mapstructure:"APP_GO_CDC_POD_HOSTNAME"`
	PodName     string `mapstructure:"APP_GO_CDC_POD_NAME"`
	PodIP       string `mapstructure:"APP_GO_CDC_POD_IP"`
	NodeName    string `mapstructure:"APP_GO_CDC_NODE_NAME"`
	Namespace   string `mapstructure:"APP_GO_CDC_POD_NAMESPACE"`

	AppLogLevel string `mapstructure:"APP_GO_CDC_LOG_LEVEL"`

	HealthCheckIntervalSeconds int `mapstructure:"APP_GO_CDC_HELTH_CHECK_INTERVAL_SECONDS"`

	AppReadEnvFromFile string `mapstructure:"APP_GO_CDC_READ_ENV_FROM_FILE"`

	DBTecnology          string `mapstructure:"APP_GO_CDC_DB_TECHNOLOGY"`
	DBHost               string `mapstructure:"APP_GO_CDC_DB_HOST"`
	DBPort               string `mapstructure:"APP_GO_CDC_DB_PORT"`
	DBUser               string `mapstructure:"APP_GO_CDC_DB_USER"`
	DBPass               string `mapstructure:"APP_GO_CDC_DB_PASS"`
	DBName               string `mapstructure:"APP_GO_CDC_DB_NAME"`
	DBMaxOpenConns       int    `mapstructure:"APP_GO_CDC_DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns       int    `mapstructure:"APP_GO_CDC_DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime    int    `mapstructure:"APP_GO_CDC_DB_CONN_MAX_LIFETIME"` // in minutes
	DBTrustServerCert    bool   `mapstructure:"APP_GO_CDC_DB_TRUST_SERVER_CERT"`
	DBConnMaxIdleTime    int    `mapstructure:"APP_GO_CDC_DB_CONN_MAX_IDLE_TIME"` // in minutes
	DBEncrypt            bool   `mapstructure:"APP_GO_CDC_DB_ENCRYPT"`
	DBPingTimeoutSeconds int    `mapstructure:"APP_GO_CDC_DB_PING_TIMEOUT_SECONDS"`
}

func getPodIP() string {
	// Fallback: detecta o IP da interface de rede principal
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "unknown"
}

// ToLogFields retorna os campos para logging estruturado
func (m *Config) ToLogFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if m.PodHostname != "" {
		fields["pod_hostname"] = m.PodHostname
	}
	if m.PodName != "" {
		fields["pod_name"] = m.PodName
	}
	if m.PodIP != "" {
		fields["pod_ip"] = m.PodIP
	}
	if m.NodeName != "" {
		fields["node_name"] = m.NodeName
	}
	if m.Namespace != "" {
		fields["namespace"] = m.Namespace
	}
	if m.AppVersion != "" {
		fields["app_version"] = m.AppVersion
	}
	if m.AppEnv != "" {
		fields["environment"] = m.AppEnv
	}

	return fields
}

// GetIdentifier retorna um identificador único do pod (para logs)
func (c *Config) GetPodIdentifier() string {
	if c.Namespace != "" {
		return c.Namespace + "/" + c.PodName
	}
	return c.PodName
}

func (c *Config) ToString(showSecrets bool) string {
	password := "******"
	if showSecrets {
		password = c.DBPass
	}

	var sb strings.Builder
	sb.WriteString("Config {\n")

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		key := ft.Tag.Get("mapstructure")
		if key == "" {
			key = ft.Name
		}

		var val string
		switch fv.Kind() {
		case reflect.String:
			val = fv.String()
			if ft.Name == "DBPass" && !showSecrets {
				val = password
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = fmt.Sprintf("%d", fv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = fmt.Sprintf("%d", fv.Uint())
		case reflect.Bool:
			val = fmt.Sprintf("%t", fv.Bool())
		default:
			val = fmt.Sprintf("%v", fv.Interface())
		}

		sb.WriteString("\t" + key + ": " + val + "\n")
	}

	sb.WriteString("}")
	return sb.String()
}

func LoadConfig(path string) (*Config, static.ErrorUtil) {
	dlog.Print("Loading config...")

	if ReadEnvFromFileEnabled() {
		dlog.Print("...from file .env")
		viper.AddConfigPath(path)
		viper.SetConfigName(".env")
		viper.SetConfigType("env")

		viper.SetDefault("APP_GO_CDC_NAME", static.APP_GO_CDC_NAME)
		viper.SetDefault("APP_GO_CDC_ENV", static.APP_GO_CDC_ENV)
		viper.SetDefault("APP_GO_CDC_LOG_LEVEL", static.APP_GO_CDC_LOG_LEVEL)
		viper.SetDefault("APP_GO_CDC_READ_ENV_FROM_FILE", "true")

		viper.SetDefault("APP_GO_CDC_HELTH_CHECK_INTERVAL_SECONDS", static.APP_GO_CDC_HELTH_CHECK_INTERVAL_SECONDS)

		viper.SetDefault("APP_GO_CDC_DB_MAX_OPEN_CONNS", static.APP_GO_CDC_DB_MAX_OPEN_CONNS)
		viper.SetDefault("APP_GO_CDC_DB_MAX_IDLE_CONNS", static.APP_GO_CDC_DB_MAX_IDLE_CONNS)
		viper.SetDefault("APP_GO_CDC_DB_CONN_MAX_LIFETIME", static.APP_GO_CDC_DB_CONN_MAX_LIFETIME)

		viper.SetDefault("APP_GO_CDC_DB_TRUST_SERVER_CERT", static.APP_GO_CDC_DB_TRUST_SERVER_CERT)
		viper.SetDefault("APP_GO_CDC_DB_ENCRYPT", static.APP_GO_CDC_DB_ENCRYPT)

		viper.SetDefault("APP_GO_CDC_DB_PING_TIMEOUT_SECONDS", static.APP_GO_CDC_DB_PING_TIMEOUT_SECONDS)
		viper.SetDefault("APP_GO_CDC_DB_CONN_MAX_IDLE_TIME", static.APP_GO_CDC_DB_CONN_MAX_IDLE_TIME)

		hostname, _ := os.Hostname()
		viper.SetDefault("APP_GO_CDC_POD_HOSTNAME", hostname)

		// Em Kubernetes, o hostname geralmente é o pod name
		podName := hostname
		viper.SetDefault("APP_GO_CDC_POD_NAME", podName)

		viper.SetDefault("APP_GO_CDC_POD_IP", getPodIP())

		viper.SetDefault("APP_GO_CDC_NODE_NAME", static.APP_GO_CDC_NODE_NAME)
		viper.SetDefault("APP_GO_CDC_POD_NAMESPACE", static.APP_GO_CDC_POD_NAMESPACE)

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				dlog.Print("Config file not found in path")
				errUtil := static.NewErrorUtil("Configuration file not found", "CONFIG_FILE_NOT_FOUND", err, err.Error())
				return nil, errUtil
			}
		}

		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			dlog.Print("Unable to decode into struct")
			errUtil := static.NewErrorUtil("Failed to decode configuration", "CONFIG_DECODE_FAILED", err, err.Error())
			return nil, errUtil
		}
		return &cfg, nil
	}

	dlog.Print("...from environment variables\n")
	var cfg Config

	appEnv := os.Getenv("APP_GO_CDC_ENV")
	if appEnv == "" {
		appEnv = static.APP_GO_CDC_ENV
	}
	cfg.AppEnv = appEnv

	cfg.AppReadEnvFromFile = "false"

	appName := os.Getenv("APP_GO_CDC_NAME")
	if appName == "" {
		appName = static.APP_GO_CDC_NAME
	}
	cfg.AppName = appName

	appLogLevel := os.Getenv("APP_GO_CDC_LOG_LEVEL")
	if appLogLevel == "" {
		appLogLevel = static.APP_GO_CDC_LOG_LEVEL
	}
	cfg.AppLogLevel = appLogLevel

	podHostname := os.Getenv("APP_GO_CDC_POD_HOSTNAME")
	if podHostname == "" {
		hostname, _ := os.Hostname()
		podHostname = hostname
	}
	cfg.PodHostname = podHostname

	podName := os.Getenv("APP_GO_CDC_POD_NAME")
	if podName == "" {
		podName = cfg.PodHostname
	}
	cfg.PodName = podName

	podIP := os.Getenv("APP_GO_CDC_POD_IP")
	if podIP == "" {
		podIP = getPodIP()
	}
	cfg.PodIP = podIP

	nodeName := os.Getenv("APP_GO_CDC_NODE_NAME")
	if nodeName == "" {
		nodeName = static.APP_GO_CDC_NODE_NAME
	}
	cfg.NodeName = nodeName

	namespace := os.Getenv("APP_GO_CDC_POD_NAMESPACE")
	if namespace == "" {
		namespace = static.APP_GO_CDC_POD_NAMESPACE
	}
	cfg.Namespace = namespace

	healthCheckInterval := os.Getenv("APP_GO_CDC_HELTH_CHECK_INTERVAL_SECONDS")
	if healthCheckInterval == "" {
		cfg.HealthCheckIntervalSeconds = static.APP_GO_CDC_HELTH_CHECK_INTERVAL_SECONDS
	} else {
		if v, err := strconv.Atoi(healthCheckInterval); err == nil {
			cfg.HealthCheckIntervalSeconds = v
		}
	}

	dbPingTimeout := os.Getenv("APP_GO_CDC_DB_PING_TIMEOUT_SECONDS")
	if dbPingTimeout == "" {
		cfg.DBPingTimeoutSeconds = static.APP_GO_CDC_DB_PING_TIMEOUT_SECONDS
	} else {
		if v, err := strconv.Atoi(dbPingTimeout); err == nil {
			cfg.DBPingTimeoutSeconds = v
		}
	}

	dbConnMaxIdleTime := os.Getenv("APP_GO_CDC_DB_CONN_MAX_IDLE_TIME")
	if dbConnMaxIdleTime == "" {
		cfg.DBConnMaxIdleTime = static.APP_GO_CDC_DB_CONN_MAX_IDLE_TIME
	} else {
		if v, err := strconv.Atoi(dbConnMaxIdleTime); err == nil {
			cfg.DBConnMaxIdleTime = v
		}
	}

	cfg.DBTecnology = os.Getenv("APP_GO_CDC_DB_TECHNOLOGY")
	cfg.DBHost = os.Getenv("APP_GO_CDC_DB_HOST")
	cfg.DBPort = os.Getenv("APP_GO_CDC_DB_PORT")
	cfg.DBUser = os.Getenv("APP_GO_CDC_DB_USER")
	cfg.DBPass = os.Getenv("APP_GO_CDC_DB_PASS")
	cfg.DBName = os.Getenv("APP_GO_CDC_DB_NAME")

	dbMaxOpenConns := os.Getenv("APP_GO_CDC_DB_MAX_OPEN_CONNS")
	if dbMaxOpenConns == "" {
		cfg.DBMaxOpenConns = static.APP_GO_CDC_DB_MAX_OPEN_CONNS
	} else {
		if v, err := strconv.Atoi(dbMaxOpenConns); err == nil {
			cfg.DBMaxOpenConns = v
		}
	}

	dbMaxIdleConns := os.Getenv("APP_GO_CDC_DB_MAX_IDLE_CONNS")
	if dbMaxIdleConns == "" {
		cfg.DBMaxIdleConns = static.APP_GO_CDC_DB_MAX_IDLE_CONNS
	} else {
		if v, err := strconv.Atoi(dbMaxIdleConns); err == nil {
			cfg.DBMaxIdleConns = v
		}
	}

	dbConnMaxLifetime := os.Getenv("APP_GO_CDC_DB_CONN_MAX_LIFETIME")
	if dbConnMaxLifetime == "" {
		cfg.DBConnMaxLifetime = static.APP_GO_CDC_DB_CONN_MAX_LIFETIME
	} else {
		if v, err := strconv.Atoi(dbConnMaxLifetime); err == nil {
			cfg.DBConnMaxLifetime = v
		}
	}

	dbTrustServerCert := os.Getenv("APP_GO_CDC_DB_TRUST_SERVER_CERT")
	if dbTrustServerCert != "" {
		if v, err := strconv.ParseBool(dbTrustServerCert); err == nil {
			cfg.DBTrustServerCert = v
		}
	}

	dbEncrypt := os.Getenv("APP_GO_CDC_DB_ENCRYPT")
	if dbEncrypt != "" {
		if v, err := strconv.ParseBool(dbEncrypt); err == nil {
			cfg.DBEncrypt = v
		}
	}

	return &cfg, nil
}
