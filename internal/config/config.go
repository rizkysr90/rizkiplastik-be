package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	AppName    string
	AppEnv     string
	ServerPort string
	LogPath    string
	PostgreSQL PostgreSQLConfig
}

// PostgreSQLConfig holds PostgreSQL database configuration
type PostgreSQLConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// GetConnectionString returns the PostgreSQL connection string
func (p PostgreSQLConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode)
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	pgMaxConns, _ := strconv.Atoi(getEnv("PG_MAX_CONNS", "10"))
	pgMinConns, _ := strconv.Atoi(getEnv("PG_MIN_CONNS", "2"))
	pgMaxConnLifetime, _ := strconv.Atoi(getEnv("PG_MAX_CONN_LIFETIME", "1800"))
	pgMaxConnIdleTime, _ := strconv.Atoi(getEnv("PG_MAX_CONN_IDLE_TIME", "30"))
	pgPort, _ := strconv.Atoi(getEnv("PG_PORT", "5432"))

	config := &Config{
		AppName:    getEnv("APP_NAME", "RizkiPlastik API"),
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogPath:    getEnv("LOG_PATH", "application.log"),
		PostgreSQL: PostgreSQLConfig{
			Host:            getEnv("PG_HOST", "localhost"),
			Port:            pgPort,
			User:            getEnv("PG_USER", "postgres"),
			Password:        getEnv("PG_PASSWORD", "postgres"),
			DBName:          getEnv("PG_DBNAME", "yourproject"),
			SSLMode:         getEnv("PG_SSLMODE", "disable"),
			MaxConns:        int32(pgMaxConns),
			MinConns:        int32(pgMinConns),
			MaxConnLifetime: time.Duration(pgMaxConnLifetime) * time.Second,
			MaxConnIdleTime: time.Duration(pgMaxConnIdleTime) * time.Second,
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
