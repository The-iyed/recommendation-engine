package graph

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
}

func (db *Neo4jDB) Recommend(ctx context.Context, base_product_id float64) ([]Product, error) {
	session := db.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
	MATCH (p:Product {product_id: $product_id})-[:SIMILAR_TO|INFO_SIMILAR_TO]-(recommended:Product)


	OPTIONAL MATCH (p)-[r1:SIMILAR_TO]-(recommended)
	OPTIONAL MATCH (p)-[r2:INFO_SIMILAR_TO]-(recommended)


	WITH recommended,
		COALESCE(SUM(r1.score), 0) AS similarity_score,
		COALESCE(SUM(r2.score), 0) AS info_similarity_score

	// Combine the scores into a total score
	WITH recommended,
		similarity_score + info_similarity_score AS total_score


	RETURN recommended, total_score
	ORDER BY total_score DESC
	LIMIT 10
`

	result, err := session.Run(query, map[string]interface{}{
		"product_id": base_product_id,
	})
	if err != nil {
		log.Println("Error running query:", err)
		return nil, err
	}

	var recommendedProducts []Product
	for result.Next() {
		record := result.Record()
		recommendedNode, _ := record.Get("recommended")
		totalScore, _ := record.Get("total_score")

		recommended := recommendedNode.(neo4j.Node)
		totalScoreValue := totalScore.(float64)
		productID := int64(recommended.Props["product_id"].(float64))

		product := Product{
			ID:          productID,
			Name:        recommended.Props["name"].(string),
			Description: recommended.Props["description"].(string),
			ImagePath:   recommended.Props["image_path"].(string),
		}

		recommendedProducts = append(recommendedProducts, product)

		log.Printf("Recommended product: %v, Score: %f", product.Name, totalScoreValue)
	}

	if err := result.Err(); err != nil {
		log.Println("Error processing query result:", err)
		return nil, err
	}

	return recommendedProducts, nil

}
