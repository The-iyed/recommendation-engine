package kafka

import (
	"context"
	"encoding/json"
	"log"
	"r-builder/modules/vector"
)

func ProcessDebeziumEvent(event KafkaMessage) {
	after := event.Payload.After

	vectorEvent := vector.VectorEvent{
		EventType:   "kafka.create.vector",
		ProductID:   after.ID,
		ImagePath:   after.ImagePath,
		Description: after.Description,
		Name:        after.Name,
		Vector:      []float64{},
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
