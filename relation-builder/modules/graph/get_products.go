package graph

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (db *Neo4jDB) GetProducts(ctx context.Context, fieldsLetter string) ([]map[string]interface{}, error) {
	session := db.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var products []map[string]interface{}

	fields := strings.Split(fieldsLetter, ".")
	var vectorFields []string
	for _, field := range fields {
		vectorFields = append(vectorFields, fmt.Sprintf("p.%s_vector", field))
	}

	additionalFieldsQuery := strings.Join(vectorFields, ", ")
	if additionalFieldsQuery != "" {
		additionalFieldsQuery = ", " + additionalFieldsQuery
	}

	query := fmt.Sprintf(`
		MATCH (p:Product)
		RETURN p.product_id, p.name, p.description, p.price, p.image_path, p.vector %s`, additionalFieldsQuery)

	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query, nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			record := result.Record()

			vectorInterface := record.Values[5].([]interface{})
			vectorFloat64 := make([]float64, len(vectorInterface))
			for i, v := range vectorInterface {
				vectorFloat64[i] = v.(float64)
			}

			product := map[string]interface{}{
				"product_id":  record.Values[0].(float64),
				"name":        record.Values[1],
				"description": record.Values[2],
				"price":       record.Values[3],
				"image_path":  record.Values[4],
				"vector":      vectorFloat64,
			}

			for i, _ := range vectorFields {
				vectorInterface := record.Values[6+i].([]interface{})
				vectorFloat64 := make([]float64, len(vectorInterface))
				for j, v := range vectorInterface {
					vectorFloat64[j] = v.(float64)
				}
				product[fields[i]+"_vector"] = vectorFloat64
			}

			products = append(products, product)
		}
		return nil, result.Err()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}
	return products, nil
}
