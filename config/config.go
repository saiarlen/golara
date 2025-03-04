package config

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Env(key string) string {
	err := godotenv.Load(".env") //loading .env file

	if err != nil {
		fmt.Print("Error loading .env file")
	}

	return os.Getenv(key)
}

// Dynamic Env Setup because .env needed rebuild if anything changes and hard to manage paas kind apps so created denv
func InitDenv() {

	viper.SetConfigName(".denv") // config file name without extension
	viper.SetConfigType("yaml")  // config file type
	viper.AddConfigPath(".")     // config file path

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigName(".denv-example")
			err = viper.ReadInConfig()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func Denv(key string) string {
	return viper.GetString(key)
}
