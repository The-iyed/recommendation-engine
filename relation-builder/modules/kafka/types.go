package kafka

import "time"

type ProductPrice struct {
	Scale int32  `json:"scale"`
	Value string `json:"value"`
}

type Product struct {
	ID          int64         `json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       *ProductPrice `json:"price"`
	ImagePath   string        `json:"image_path"`
}

type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	DB        string `json:"db"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	TxID      *int64 `json:"txId"`
	Lsn       *int64 `json:"lsn"`
	Xmin      *int64 `json:"xmin"`
}

type Transaction struct {
	ID                  string `json:"id"`
	TotalOrder          int64  `json:"total_order"`
	DataCollectionOrder int64  `json:"data_collection_order"`
}

type Envelope struct {
	Before      *Product     `json:"before"`
	After       *Product     `json:"after"`
	Source      Source       `json:"source"`
	Op          string       `json:"op"`
	TsMs        *int64       `json:"ts_ms"`
	Transaction *Transaction `json:"transaction"`
}

type SchemaField struct {
	Field string `json:"field"`
}

type Schema struct {
	Fields []SchemaField `json:"fields"`
}

type Payload struct {
	Before      *Product     `json:"before"`
	After       *Product     `json:"after"`
	Op          string       `json:"op"`
	Source      Source       `json:"source"`
	TsMs        int64        `json:"ts_ms"`
	Transaction *Transaction `json:"transaction"`
}

type KafkaMessage struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

type VectorData struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImagePath   string    `json:"image_path"`
	Vector      []float64 `json:"vector"`
}
