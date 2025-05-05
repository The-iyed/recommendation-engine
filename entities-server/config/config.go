package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	Env                   string
	DB_HOST               string
	DB_PORT               string
	DB_PASSWORD           string
	DB_NAME               string
	DB_USER               string
	PATH                  string
	CLOUDINARY_CLOUD_NAME string
	CLOUDINARY_API_KEY    string
	CLOUDINARY_API_SECRET string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		Port:                  getEnv("PORT", "5002"),
		Env:                   getEnv("ENV", "dev"),
		DB_PORT:               getEnv("DB_PORT", "5432"),
		DB_HOST:               getEnv("DB_HOST", "postgres"),
		DB_USER:               getEnv("DB_USER", "postgres"),
		DB_NAME:               getEnv("DB_NAME", "productdb"),
		DB_PASSWORD:           getEnv("DB_PASSWORD", "root1234"),
		PATH:                  getEnv("PATH", "http://localhost:5001/"),
		CLOUDINARY_CLOUD_NAME: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CLOUDINARY_API_KEY:    getEnv("CLOUDINARY_API_KEY", ""),
		CLOUDINARY_API_SECRET: getEnv("CLOUDINARY_API_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
