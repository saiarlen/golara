package config

import (
	"github.com/test/myapp/internal/constants"
	"os"
	"strconv"
	"strings"
)

// AppConfig holds application configuration
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	URL         string
	Port        string
	Timezone    string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Default     string
	Connections map[string]DatabaseConnection
}

type DatabaseConnection struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Charset  string
	Options  map[string]string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Default string
	Stores  map[string]CacheStore
}

type CacheStore struct {
	Driver   string
	Host     string
	Port     string
	Password string
	Database int
	Prefix   string
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	Default     string
	Connections map[string]QueueConnection
}

type QueueConnection struct {
	Driver   string
	Host     string
	Port     string
	Password string
	Database int
	Queue    string
}

// LoadAppConfig loads application configuration
func LoadAppConfig() *AppConfig {
	return &AppConfig{
		Name:        GetEnv("APP_NAME", "Golara Application"),
		Environment: GetEnv("APP_ENV", constants.EnvDevelopment),
		Debug:       GetEnvBool("APP_DEBUG", true),
		URL:         GetEnv("APP_URL", "http://localhost:9000"),
		Port:        GetEnv("APP_PORT", constants.DefaultPort),
		Timezone:    GetEnv("APP_TIMEZONE", "UTC"),
	}
}

// LoadDatabaseConfig loads database configuration
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Default: GetEnv("DB_CONNECTION", "mysql"),
		Connections: map[string]DatabaseConnection{
			"mysql": {
				Driver:   "mysql",
				Host:     GetEnv("DB_HOST", "localhost"),
				Port:     GetEnv("DB_PORT", "3306"),
				Database: GetEnv("DB_DATABASE", "golara"),
				Username: GetEnv("DB_USERNAME", "root"),
				Password: GetEnv("DB_PASSWORD", ""),
				Charset:  GetEnv("DB_CHARSET", "utf8mb4"),
				Options: map[string]string{
					"parseTime": "True",
					"loc":       "Local",
				},
			},
		},
	}
}

// LoadCacheConfig loads cache configuration
func LoadCacheConfig() *CacheConfig {
	return &CacheConfig{
		Default: GetEnv("CACHE_DRIVER", constants.CacheDriverMemory),
		Stores: map[string]CacheStore{
			constants.CacheDriverRedis: {
				Driver:   constants.CacheDriverRedis,
				Host:     GetEnv("REDIS_HOST", "localhost"),
				Port:     GetEnv("REDIS_PORT", "6379"),
				Password: GetEnv("REDIS_PASSWORD", ""),
				Database: GetEnvInt("REDIS_DB", 0),
				Prefix:   GetEnv("CACHE_PREFIX", constants.DefaultCachePrefix),
			},
			constants.CacheDriverMemory: {
				Driver: constants.CacheDriverMemory,
				Prefix: GetEnv("CACHE_PREFIX", constants.DefaultCachePrefix),
			},
		},
	}
}

// LoadQueueConfig loads queue configuration
func LoadQueueConfig() *QueueConfig {
	return &QueueConfig{
		Default: GetEnv("QUEUE_CONNECTION", constants.QueueDriverMemory),
		Connections: map[string]QueueConnection{
			constants.QueueDriverRedis: {
				Driver:   constants.QueueDriverRedis,
				Host:     GetEnv("REDIS_HOST", "localhost"),
				Port:     GetEnv("REDIS_PORT", "6379"),
				Password: GetEnv("REDIS_PASSWORD", ""),
				Database: GetEnvInt("REDIS_DB", 0),
				Queue:    GetEnv("QUEUE_NAME", constants.DefaultQueueName),
			},
		},
	}
}

// Helper functions for environment variables
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := Denv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if value := GetEnv(key, ""); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := GetEnv(key, ""); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func GetEnvSlice(key string, defaultValue []string, separator string) []string {
	if value := GetEnv(key, ""); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}

// IsProduction checks if app is in production
func IsProduction() bool {
	return GetEnv("APP_ENV", constants.EnvDevelopment) == constants.EnvProduction
}

// IsDevelopment checks if app is in development
func IsDevelopment() bool {
	return GetEnv("APP_ENV", constants.EnvDevelopment) == constants.EnvDevelopment
}

// IsTesting checks if app is in testing
func IsTesting() bool {
	return GetEnv("APP_ENV", constants.EnvDevelopment) == constants.EnvTesting
}
