package kafka

import (
	"context"
	"log"
	"r-builder/config"
	"r-builder/modules/graph"

	"github.com/segmentio/kafka-go"
)

type IKafkaConsumer struct {
	Reader *kafka.Reader
	Neo4j  *graph.Neo4jDB
}

func KafkaConsumer() *IKafkaConsumer {
	cfg := config.LoadConfig()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroup,
		Topic:   cfg.KafkaTopic,
	})
	neo4j, err := graph.NewNeo4jDB()
	if err != nil {
		log.Fatal("Failed to establish connection with neo4J es5aaaat :", err)

	}
	log.Println("Kafka consumer is connected and ready to consume messages...")

	return &IKafkaConsumer{
		Reader: reader,
		Neo4j:  neo4j,
	}
}

func (k *IKafkaConsumer) ConsumeMessages(ctx context.Context, handleEntityCDC func(KafkaMessage)) {
	for {
		msg, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while consuming message: %v", err)
			continue
		}
		log.Printf("Consumed message from topic %s:", msg.Topic)

		switch msg.Topic {
		case "kafka.vector.created":
			if err := k.processVectorCreatedMessage(ctx, msg); err != nil {
				log.Printf("[Kafka Consumer] Error processing vector created message (offset: %d): %v", msg.Offset, err)
			}

		}
	}
}

func (k *IKafkaConsumer) Close() {
	err := k.Reader.Close()
	if err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
	}
	log.Println("Kafka consumer closed.")
}
