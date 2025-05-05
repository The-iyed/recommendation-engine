package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	Redis_Host   string
	Redis_Port   string
	KafkaBroker  string
	KafkaGroup   string
	FieldsLetter string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		Port:         getEnv("PORT", "5002"),
		Redis_Host:   getEnv("REDIS_HOST", "redis"),
		Env:          getEnv("ENV", "dev"),
		Redis_Port:   getEnv("REDIS_PORT", "6379"),
		KafkaBroker:  getEnv("KAFKA_BROKER", "kafka:29092"),
		KafkaGroup:   getEnv("KAFKA_GROUPE", "recommendation-server"),
		FieldsLetter: getEnv("FIELDS_LETTER", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
