package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"r-builder/config"
)

type IKafkaProducer struct {
	Writer *kafka.Writer
}

func KafkaProducer() *IKafkaProducer {
	cfg := config.LoadConfig()
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Balancer: &kafka.LeastBytes{},
	})

	log.Println("Kafka producer is connected and ready to produce messages...")

	return &IKafkaProducer{
		Writer: writer,
	}
}

func (k *IKafkaProducer) ProduceMessage(ctx context.Context, topic string, message []byte) error {
	err := k.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: message,
	})

	if err != nil {
		log.Printf("error while producing message: %v", err)
		return err
	}

	log.Printf("Message sent to topic %s", topic)
	return nil
}

func (k *IKafkaProducer) Close() {
	k.Writer.Close()
}
