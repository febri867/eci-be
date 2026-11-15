package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	MongoURI           string
	MongoDatabase      string
	TurnstileSecretKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return &Config{
		Port:               getEnv("PORT", "4000"),
		MongoURI:           getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase:      getEnv("MONGODB_DATABASE", "eci"),
		TurnstileSecretKey: getEnv("TURNSTILE_SECRET_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
