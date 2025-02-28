package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBURL      string
	GroqAPIKey string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AppConfig = Config{
		Port:       viper.GetString("PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBURL:      viper.GetString("DB_URL"),
		GroqAPIKey: viper.GetString("GROQ_API_KEY"),
	}
}
