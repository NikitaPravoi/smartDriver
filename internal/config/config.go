package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration settings
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Centrifugo CentrifugoConfig
	Parser     ParserConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Environment  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CentrifugoConfig struct {
	URL    string
	APIKey string
}

type ParserConfig struct {
	PollingInterval time.Duration
	BatchSize       int
	WorkerCount     int
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{}
	var err error

	// Server configuration
	cfg.Server = ServerConfig{
		Port:         getEnvOrDefault("SERVER_PORT", "8080"),
		Host:         getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
		ReadTimeout:  getDurationOrDefault("SERVER_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationOrDefault("SERVER_WRITE_TIMEOUT", 15*time.Second),
		Environment:  getEnvOrDefault("APP_ENV", "development"),
	}

	// Database configuration
	cfg.Database = DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getRequiredEnv("DB_USER"),
		Password: getRequiredEnv("DB_PASSWORD"),
		DBName:   getRequiredEnv("DB_NAME"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	// Centrifugo configuration
	cfg.Centrifugo = CentrifugoConfig{
		URL:    getRequiredEnv("CENTRIFUGO_URL"),
		APIKey: getRequiredEnv("CENTRIFUGO_API_KEY"),
	}

	// Parser configuration
	cfg.Parser = ParserConfig{
		PollingInterval: getDurationOrDefault("PARSER_POLLING_INTERVAL", 30*time.Second),
		BatchSize:       getIntOrDefault("PARSER_BATCH_SIZE", 100),
		WorkerCount:     getIntOrDefault("PARSER_WORKER_COUNT", 5),
	}

	return cfg, err
}

// Helper functions for environment variables
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntOrDefault(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}

func getDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDSN returns the database connection string
func (c DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
