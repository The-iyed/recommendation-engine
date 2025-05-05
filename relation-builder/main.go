package main

import (
	"context"
	"log"
	"r-builder/config"
	"r-builder/modules/kafka"
	"r-builder/modules/redis"
	"r-builder/router"
)

func main() {
	cfg := config.LoadConfig()

	redis.InitRedis()
	ctx := context.Background()

	kafkaConsumer := kafka.KafkaConsumer()
	kafka.InitializeTopics()
	go func() {
		kafkaConsumer.ConsumeMessages(ctx, kafka.ProcessDebeziumEvent)

	}()
	r := router.SetupRouter()
	log.Printf("Starting server on port %s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
