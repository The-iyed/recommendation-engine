package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Env           string
	Redis_Host    string
	Redis_Port    string
	KafkaBroker   string
	KafkaGroup    string
	KafkaTopic    string
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		Port:          getEnv("PORT", "5004"),
		Redis_Host:    getEnv("REDIS_HOST", "redis"),
		Env:           getEnv("ENV", "dev"),
		Redis_Port:    getEnv("REDIS_PORT", "6379"),
		KafkaBroker:   getEnv("KAFKA_BROKER", "kafka:29092"),
		KafkaGroup:    getEnv("KAFKA_GROUPE", "relation-builder"),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "kafka.vector.created"),
		Neo4jURI:      getEnv("Neo4jURI", ""),
		Neo4jUser:     getEnv("Neo4jUser", ""),
		Neo4jPassword: getEnv("Neo4jPassword", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
