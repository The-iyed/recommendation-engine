package graph

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (n *Neo4jDB) StoreProduct(ctx context.Context, event map[string]interface{}) error {

	session := n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	tx, err := session.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Close()

	params := map[string]interface{}{
		"product_id":  event["id"],
		"name":        event["name"],
		"description": event["description"],
		"price":       event["price"],
		"image_path":  event["image_path"],
		"vector":      event["vector"],
		"fields_letter": event["fields_letter"],
	}

	dynamicFields := ""
	for key, value := range event {
		if strings.HasSuffix(key, "_vector") {
			fieldName := strings.ReplaceAll(key, ".", "_") 
			dynamicFields += fmt.Sprintf(", p.%s = $%s", fieldName, fieldName)
			params[fieldName] = value
		}
	}

	query := fmt.Sprintf(`
		MERGE (p:Product {product_id: $product_id})
		SET p.name = $name,
			p.description = $description,
			p.price = $price,
			p.image_path = $image_path,
			p.fields_letter = $fields_letter,
			p.vector = $vector
			%s`, dynamicFields)

	_, err = tx.Run(query, params)
	if err != nil {
		return fmt.Errorf("error storing product and vectors in Neo4j: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Printf("Successfully stored product and vectors for product_id %s", event["id"])
	return nil
}
