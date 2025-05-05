package main

import (
	"context"
	"log"
	"r-server/config"
	"r-server/modules/kafka"
	"r-server/modules/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()
	redis.InitRedis()
	ctx := context.Background()

	kafkaConsumer := kafka.KafkaConsumer()
	kafka.InitializeTopics()
	go func() {
		kafkaConsumer.ConsumeMessages(ctx, kafka.ProcessDebeziumEvent)
		
	}()

	log.Printf("Starting server on port %s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
