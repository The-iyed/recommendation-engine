package vector

type VectorEvent struct {
	EventType   string    `json:"event_type"`
	ProductID   int64     `json:"product_id"`
	ImagePath   string    `json:"image_path"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Vector      []float64 `json:"vector"`
}
