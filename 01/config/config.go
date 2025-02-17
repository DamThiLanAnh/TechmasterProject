package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, loading default values")
	}
}
