package entity

type Product struct {
	Sku           string  `json:"sku" database:"sku"`
	Name          string  `json:"name" database:"name"`
	Description   string  `json:"description" database:"description"`
	Price         float64 `json:"price" database:"price"`
	StockQuantity int     `json:"stock_quantity" database:"stock_quantity"`
	CreatedAt     string  `json:"created_at" database:"created_at"`
	UpdatedAt     string  `json:"updated_at" database:"updated_at"`
}
