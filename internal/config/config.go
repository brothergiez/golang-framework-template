package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ServiceTypeAPI    = "api"
	ServiceTypeCron   = "cron"
	ServiceTypeWorker = "worker"
	ServiceTypeLambda = "lambda"
)

type Config struct {
	App      AppConfig
	Service  ServiceConfig
	HTTP     HTTPConfig
	Cron     CronConfig
	Worker   WorkerConfig
	Lambda   LambdaConfig
	Database DatabaseConfig
	Mongo    MongoConfig
	Redis    RedisConfig
	Logger   LoggerConfig
}

type AppConfig struct {
	Name            string
	Env             string
	ShutdownTimeout time.Duration
}

type ServiceConfig struct {
	Type string
}

type HTTPConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BodyLimit    int64
}

type CronConfig struct {
	Timezone string
}

type WorkerConfig struct {
	Concurrency int
	PollInterval time.Duration
}

type LambdaConfig struct {
	Handler string
}

type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type MongoConfig struct {
	Enabled  bool
	URI      string
	Database string
}

type RedisConfig struct {
	Enabled  bool
	URL      string
	Database int
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "go-framework"),
			Env:             getEnv("APP_ENV", "development"),
			ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Service: ServiceConfig{
			Type: strings.ToLower(getEnv("SERVICE_TYPE", ServiceTypeAPI)),
		},
		HTTP: HTTPConfig{
			Host:         getEnv("HTTP_HOST", "0.0.0.0"),
			Port:         getInt("APP_PORT", 8080),
			ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			BodyLimit:    getInt64("HTTP_BODY_LIMIT", 1<<20),
		},
		Cron: CronConfig{
			Timezone: getEnv("CRON_TIMEZONE", "UTC"),
		},
		Worker: WorkerConfig{
			Concurrency:  getInt("WORKER_CONCURRENCY", 4),
			PollInterval: getDuration("WORKER_POLL_INTERVAL", 2*time.Second),
		},
		Lambda: LambdaConfig{
			Handler: getEnv("LAMBDA_HANDLER", "default"),
		},
		Database: DatabaseConfig{
			Driver:          strings.ToLower(getEnv("DATABASE_DRIVER", "memory")),
			Host:            getEnv("DATABASE_HOST", "localhost"),
			Port:            getInt("DATABASE_PORT", 5432),
			Name:            getEnv("DATABASE_NAME", "mydb"),
			User:            getEnv("DATABASE_USER", "postgres"),
			Password:        getEnv("DATABASE_PASSWORD", "password"),
			SSLMode:         getEnv("DATABASE_SSLMODE", "disable"),
			MaxOpenConns:    getInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DATABASE_CONN_MAX_LIFETIME", time.Hour),
			ConnMaxIdleTime: getDuration("DATABASE_CONN_MAX_IDLE_TIME", 30*time.Minute),
		},
		Mongo: MongoConfig{
			Enabled:  getBool("MONGO_ENABLED", false),
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "mydb"),
		},
		Redis: RedisConfig{
			Enabled:  getBool("REDIS_ENABLED", false),
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Database: getInt("REDIS_DATABASE", 0),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	switch c.Service.Type {
	case ServiceTypeAPI, ServiceTypeCron, ServiceTypeWorker, ServiceTypeLambda:
	default:
		return fmt.Errorf("invalid SERVICE_TYPE %q: must be api|cron|worker|lambda", c.Service.Type)
	}

	switch c.Database.Driver {
	case "memory", "postgres", "mysql", "mariadb", "mssql", "mongo":
	default:
		return fmt.Errorf("invalid DATABASE_DRIVER %q", c.Database.Driver)
	}

	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if c.Service.Type == ServiceTypeAPI && (c.HTTP.Port < 1 || c.HTTP.Port > 65535) {
		return fmt.Errorf("APP_PORT must be between 1 and 65535")
	}
	return nil
}

func (c *Config) HTTPAddr() string {
	return fmt.Sprintf("%s:%d", c.HTTP.Host, c.HTTP.Port)
}

func (d DatabaseConfig) DSN() string {
	switch d.Driver {
	case "postgres":
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
		)
	case "mysql", "mariadb":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
			d.User, d.Password, d.Host, d.Port, d.Name,
		)
	case "mssql":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			d.User, d.Password, d.Host, d.Port, d.Name,
		)
	default:
		return ""
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return n
}

func getBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
