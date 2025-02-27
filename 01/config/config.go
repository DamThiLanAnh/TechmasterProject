package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config lưu trữ các biến cấu hình
type Config struct {
	GroqAPIKey string
}

// LoadConfig đọc cấu hình từ file .env
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Không tìm thấy file .env, sẽ dùng biến môi trường")
	}

	apiKey := viper.GetString("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY chưa được thiết lập")
	}

	return &Config{GroqAPIKey: apiKey}
}
