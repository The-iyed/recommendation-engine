package services

import (
	"context"
	"log"
	"r-builder/modules/graph"
)

func RecommendProducts(ctx context.Context, db *graph.Neo4jDB, baseProductID float64) ([]graph.Product, error) {

	products, err := db.Recommend(ctx, baseProductID)
	if err != nil {
		log.Println("Error fetching recommended products:", err)
		return nil, err
	}

	return products, nil
}
