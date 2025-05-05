package graph

import (
	"context"
	"fmt"
	"log"
	"r-builder/pkg"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Tempo struct {
	ProductID int64 `json:"product_id"`
}

func (db *Neo4jDB) BuildSimilarities(ctx context.Context, products []map[string]interface{}, new_product map[string]interface{}, w_info float64, w_image float64, fieldsLetter string) error {
	session := db.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	fields := strings.Split(fieldsLetter, ".")

	for _, product := range products {

		v1, err := pkg.InterfaceToFloat64Slice(new_product["vector"])
		if err != nil {
			log.Println("Error while parsing new product vector:", err)
			return err
		}
		v2, err := pkg.InterfaceToFloat64Slice(product["vector"])
		if err != nil {
			log.Println("Error while parsing product vector:", err)
			return err
		}
		if len(v1) == 0 || len(v2) == 0 {
			log.Println("Empty vector encountered, skipping similarity calculation")
			continue
		}

		imageSimilarity := pkg.CosineSimilarity(v1, v2)

		var totalSimilarity float64
		for _, field := range fields {

			newFieldVector, _ := pkg.InterfaceToFloat64Slice(new_product[field+"_vector"])
			productFieldVector, _ := pkg.InterfaceToFloat64Slice(product[field+"_vector"])
			if len(newFieldVector) == 0 || len(productFieldVector) == 0 {
				log.Printf("Empty vector for field %s, skipping similarity calculation", field)
				continue
			}
			totalSimilarity += pkg.CosineSimilarity(newFieldVector, productFieldVector)
		}

		infoSimilarity := totalSimilarity / float64(len(fields))

		id1, ok1 := convertToInt64(new_product["product_id"])
		id2, ok2 := convertToInt64(product["product_id"])
		if ok1 != nil || ok2 != nil {
			log.Println("Invalid product_id format, expected int64")
			continue
		}
		if imageSimilarity >= w_image && new_product["product_id"] != product["product_id"] {
			_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

				query := `
						MATCH (p1:Product {product_id: $id1}), (p2:Product {product_id: $id2})
						MERGE (p1)-[r1:SIMILAR_TO {score: $score}]->(p2)
						MERGE (p2)-[r2:SIMILAR_TO {score: $score}]->(p1)
				`
				_, err := tx.Run(query, map[string]interface{}{
					"id1":   id1,
					"id2":   id2,
					"score": imageSimilarity,
				})
				return nil, err
			})
			if err != nil {
				log.Printf("Error storing IMAGE_SIMILAR_TO relation for product %v and %v: %v", new_product["product_id"], product["product_id"], err)
			}
		}

		if infoSimilarity >= w_info && new_product["product_id"] != product["product_id"] {
			_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				query := `
						MATCH (p1:Product {product_id: $id1}), (p2:Product {product_id: $id2})
						MERGE (p1)-[r1:INFO_SIMILAR_TO {score: $score}]->(p2)
						MERGE (p2)-[r2:INFO_SIMILAR_TO {score: $score}]->(p1)

				`
				_, err := tx.Run(query,
					map[string]interface{}{
						"id1":   id1,
						"id2":   id2,
						"score": infoSimilarity,
					})
				return nil, err
			})
			if err != nil {
				log.Printf("Error storing INFO_SIMILAR_TO relation for product %v and %v: %v", new_product["product_id"], product["product_id"], err)
			}
		}
	}

	return nil
}

func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported product_id type: %T", value)
	}
}
