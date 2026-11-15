package entity

type OrderItemEntity struct {
	OrderId   string `json:"order_id" database:"order_id"`
	Sku       string `json:"sku" database:"sku"`
	Quantity  int    `json:"quantity" database:"quantity"`
	CreatedAt string `json:"created_at" database:"created_at"`
	UpdatedAt string `json:"updated_at" database:"updated_at"`
}
