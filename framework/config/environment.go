package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// EnvironmentManager manages multi-environment configuration
type EnvironmentManager struct {
	environment string
	configs     map[string]*viper.Viper
	basePath    string
}

// NewEnvironmentManager creates a new environment manager
func NewEnvironmentManager(basePath string) *EnvironmentManager {
	return &EnvironmentManager{
		configs:  make(map[string]*viper.Viper),
		basePath: basePath,
	}
}

// LoadEnvironment loads configuration for specific environment
func (em *EnvironmentManager) LoadEnvironment(env string) error {
	em.environment = env
	
	// Load base configuration
	if err := em.loadConfig("config", ""); err != nil {
		return fmt.Errorf("failed to load base config: %w", err)
	}
	
	// Load environment-specific configuration
	if env != "" {
		if err := em.loadConfig("config", env); err != nil {
			log.Printf("Warning: No environment-specific config found for %s", env)
		}
	}
	
	// Load additional config files
	configFiles := []string{"database", "cache", "queue", "mail", "services"}
	for _, configFile := range configFiles {
		em.loadConfig(configFile, "")
		if env != "" {
			em.loadConfig(configFile, env)
		}
	}
	
	log.Printf("✅ Configuration loaded for environment: %s", env)
	return nil
}

func (em *EnvironmentManager) loadConfig(name, env string) error {
	v := viper.New()
	
	// Set config name and type
	if env != "" {
		v.SetConfigName(fmt.Sprintf("%s.%s", name, env))
	} else {
		v.SetConfigName(name)
	}
	
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(em.basePath, "config"))
	v.AddConfigPath(".")
	
	// Try to read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return err // Config file not found
		}
		return fmt.Errorf("error reading config file: %w", err)
	}
	
	// Store config
	configKey := name
	if env != "" {
		configKey = fmt.Sprintf("%s.%s", name, env)
	}
	em.configs[configKey] = v
	
	return nil
}

// Get retrieves a configuration value
func (em *EnvironmentManager) Get(key string) interface{} {
	// Try environment-specific config first
	if em.environment != "" {
		parts := strings.Split(key, ".")
		if len(parts) > 0 {
			envConfigKey := fmt.Sprintf("%s.%s", parts[0], em.environment)
			if config, exists := em.configs[envConfigKey]; exists {
				if value := config.Get(strings.Join(parts[1:], ".")); value != nil {
					return value
				}
			}
		}
	}
	
	// Fall back to base config
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		if config, exists := em.configs[parts[0]]; exists {
			return config.Get(strings.Join(parts[1:], "."))
		}
	}
	
	return nil
}

// GetString retrieves a string configuration value
func (em *EnvironmentManager) GetString(key string) string {
	if value := em.Get(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// DetectEnvironment detects environment from various sources
func DetectEnvironment() string {
	// Check environment variable
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		return env
	}
	
	// Default to development
	return "development"
}