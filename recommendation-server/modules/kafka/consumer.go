package kafka

import (
	"context"
	"encoding/json"
	"log"
	"r-server/config"

	"github.com/segmentio/kafka-go"
)

type IKafkaConsumer struct {
	Reader *kafka.Reader
}

func KafkaConsumer() *IKafkaConsumer {
	cfg := config.LoadConfig()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		GroupID: cfg.KafkaGroup,
		Topic:   "productdb_server.public.products",
	})

	log.Println("Kafka consumer is connected and ready to consume messages...")

	return &IKafkaConsumer{
		Reader: reader,
	}
}

func (k *IKafkaConsumer) ConsumeMessages(ctx context.Context, handleEntityCDC func(KafkaMessage)) {
	for {
		msg, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while consuming message: %v", err)
			continue
		}

		switch msg.Topic {
		case "productdb_server.public.products":
			var event KafkaMessage
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("error unmarshalling message: %v | Message content: %s", err, string(msg.Value))
				continue
			}
			handleEntityCDC(event)
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

func ProcessDebeziumEvent(event KafkaMessage) {
	after := event.Payload.After
	cfg := config.LoadConfig()
	vectorEvent := VectorEvent{
		EventType:    "kafka.create.vector",
		ProductID:    after.ID,
		ImagePath:    after.ImagePath,
		Description:  after.Description,
		Name:         after.Name,
		Vector:       "empty_vector",
		FieldsLetter: cfg.FieldsLetter,
	}

	message, err := json.Marshal(vectorEvent)
	if err != nil {
		log.Printf("Error marshalling vector event: %v", err)
		return
	}

	producer := KafkaProducer()
	defer producer.Close()

	err = producer.ProduceMessage(context.Background(), "kafka.create.vector", message)
	if err != nil {
		log.Printf("Error producing message to Kafka: %v", err)
	}
}
