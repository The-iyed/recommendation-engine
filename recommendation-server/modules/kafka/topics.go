package kafka

import (
	"fmt"
	"log"
	"r-server/config"

	"github.com/segmentio/kafka-go"
)

func InitializeTopics() {
	cfg := config.LoadConfig()
	conn, err := kafka.Dial("tcp", cfg.KafkaBroker)
	if err != nil {
		log.Fatal("Failed to dial kafka :", err)
	}
	defer conn.Close()
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             "productdb_server.public.products",
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
		{
			Topic:             "kafka.create.vector",
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}
	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatal("Failed to create topics :", err)
	}
	fmt.Println("Topics created successfully !")
}
