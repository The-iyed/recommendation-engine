# Recommendation System

A comprehensive microservices-based recommendation engine that uses image vectorization and textual similarity to provide accurate product recommendations. The system employs Change Data Capture (CDC), vector embeddings, and graph relationships to create a robust recommendation platform.

## System Architecture

The recommendation system consists of several microservices working together:

1. **Entities Server**: Manages product data in PostgreSQL
2. **Vectorization Server**: Converts product images and text fields into vector embeddings
3. **Relation Builder**: Creates and manages similarity relationships between products
4. **Recommendation Server**: Handles recommendation requests and processing
5. **Infrastructure**: Contains shared services like Kafka, Zookeeper, Redis, PostgreSQL, and Neo4j

## Core Recommendation Features

### Real-time Data Processing with CDC (Change Data Capture)

The system uses Debezium to capture changes in the PostgreSQL database in real-time. When a product is added, updated, or deleted:

1. Debezium captures the change and sends it to Kafka
2. The Recommendation Server consumes these events from Kafka
3. The data flows through the recommendation pipeline for processing

This CDC approach ensures recommendations remain current without requiring batch processing or manual triggers.

### Multi-factor Recommendation Algorithm

The system recommends products based on two primary factors:

1. **Image Similarity**: Uses ResNet50 to generate vector embeddings from product images
2. **Text Field Similarity**: Uses Sentence Transformers to create vector embeddings from configurable text fields (e.g., name, description)

The fields used for similarity calculations are configurable in the `.env` file through the `FIELDS_LETTER` parameter. For example, `FIELDS_LETTER=name.description` configures the system to use both name and description fields for text similarity.

### Similarity Calculation and Weighting

The system calculates similarity using cosine similarity between vector embeddings:

1. **Cosine Similarity**: Measures the cosine of the angle between two vectors, providing a similarity score between -1 and 1
   ```go
   func CosineSimilarity(vec1, vec2 []float64) float64 {
       var normAB, normA, normB float64
       for i := 0; i < len(vec1); i++ {
           normAB += vec1[i] * vec2[i]
           normA += vec1[i] * vec1[i]
           normB += vec2[i] * vec2[i]
       }
       if normA == 0 || normB == 0 {
           return 0
       }
       return normAB / (math.Sqrt(normA) * math.Sqrt(normB))
   }
   ```

2. **Configurable Weights**: The system allows setting thresholds for image and text similarity to determine when relationships should be created
   - Image similarity threshold (`w_image`)
   - Text information similarity threshold (`w_info`)

### Graph-based Relationship Management

The system uses Neo4j graph database to store and query relationships between products:

1. **Relationship Types**:
   - `SIMILAR_TO`: Represents image-based similarity
   - `INFO_SIMILAR_TO`: Represents text-based similarity

2. **Recommendation Query**:
   The recommendation algorithm combines both similarity scores to provide a total score for recommendations:
   ```cypher
   MATCH (p:Product {product_id: $product_id})-[:SIMILAR_TO|INFO_SIMILAR_TO]-(recommended:Product)
   OPTIONAL MATCH (p)-[r1:SIMILAR_TO]-(recommended)
   OPTIONAL MATCH (p)-[r2:INFO_SIMILAR_TO]-(recommended)
   WITH recommended,
       COALESCE(SUM(r1.score), 0) AS similarity_score,
       COALESCE(SUM(r2.score), 0) AS info_similarity_score
   WITH recommended,
       similarity_score + info_similarity_score AS total_score
   RETURN recommended, total_score
   ORDER BY total_score DESC
   LIMIT 10
   ```

## Data Flow Pipeline

1. **Product Creation/Update**:
   - Product data is stored in PostgreSQL
   - Debezium captures the change and publishes to Kafka topic `productdb_server.public.products`

2. **Vectorization**:
   - Recommendation Server consumes the product event and publishes to `kafka.create.vector` topic
   - Vectorization Server processes the product data:
     - Generates image vectors using ResNet50
     - Generates text vectors for specified fields (name, description) using Sentence Transformers
     - Publishes vector data to `kafka.vector.created` topic

3. **Relationship Building**:
   - Relation Builder consumes the vector data
   - Calculates similarity with existing products
   - Creates relationships in Neo4j when similarity exceeds thresholds
   - Stores recommendations in Redis for quick access

4. **Recommendation Serving**:
   - When a recommendation is requested for a product, the system queries Neo4j
   - Returns products with highest combined similarity scores

## Setup and Configuration

The system uses environment variables for configuration:

- **Field Selection**: `FIELDS_LETTER` determines which text fields are used for similarity calculation
- **Vectorization Models**: 
  - Image vectors: ResNet50
  - Text vectors: Configurable (default: all-MiniLM-L6-v2)
- **Vector Dimensions**: Configurable vector sizes
- **Similarity Thresholds**: Determine when to create relationships between products

## Getting Started

Run the system using the provided start script:

```bash
./start.sh
```

This will start all required services in Docker containers.

## Current Status

This project is currently in its core development phase, focusing on the fundamental recommendation engine functionality. Future enhancements may include:
- User preference integration
- A/B testing framework
- Performance optimization
- Additional recommendation algorithms
- User interface for system configuration
