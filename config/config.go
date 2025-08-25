package config

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	envOnce sync.Once
	envErr  error
)

// Env loads .env file once and returns environment variable
func Env(key string) string {
	envOnce.Do(func() {
		envErr = godotenv.Load(".env")
		if envErr != nil {
			log.Printf("Warning: Error loading .env file: %v", envErr)
		}
	})

	return os.Getenv(key)
}

// EnvWithDefault returns environment variable with default value
func EnvWithDefault(key, defaultValue string) string {
	value := Env(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// InitDenv initializes dynamic environment configuration
func InitDenv() error {
	viper.SetConfigName(".denv")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigName(".denv-example")
			err = viper.ReadInConfig()
			if err != nil {
				return err
			}
			log.Println("Using .denv-example.yaml as configuration")
		} else {
			return err
		}
	} else {
		log.Println("Configuration loaded from .denv.yaml")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuration file changed: %s", e.Name)
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reloading config: %v", err)
		}
	})
	
	return nil
}

func Denv(key string) string {
	return viper.GetString(key)
}