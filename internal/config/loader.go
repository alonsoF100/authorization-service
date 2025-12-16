package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found", err)
	}

	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Unable to decode config:", err)
	}

	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")

	log.Println("Config loaded successfully")
	return &config
}
