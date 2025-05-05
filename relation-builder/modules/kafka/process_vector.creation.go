package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func (k *IKafkaConsumer) processVectorCreated(ctx context.Context, event map[string]interface{}) error {
	if err := k.Neo4j.StoreProduct(ctx, event); err != nil {
		return err
	}
	fields_letter := event["fields_letter"].(string)
	log.Println("product stored")
	products, err := k.Neo4j.GetProducts(ctx, fields_letter)
	if err != nil {
		return err
	}

	err = k.Neo4j.BuildSimilarities(ctx, products, event, 0.4, 0.7, fields_letter)
	if err != nil {
		return err
	}
	log.Println("product similiraity built")

	return nil
}

func (k *IKafkaConsumer) processVectorCreatedMessage(ctx context.Context, msg kafka.Message) error {
	var vectorData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &vectorData); err != nil {
		return err
	}

	return k.processVectorCreated(ctx, vectorData)
}
