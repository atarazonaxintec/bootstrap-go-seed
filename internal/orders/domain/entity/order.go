package entity

type OrderEntity struct {
	Id         string               `json:"id" database:"id"`
	Rut        string               `json:"rut" database:"rut"`
	SubTotal   float64              `json:"sub_total" database:"sub_total"`
	Status     OrderStatus          `json:"status" database:"status"`
	CreatedAt  string               `json:"created_at" database:"created_at"`
	UpdatedAt  string               `json:"updated_at" database:"updated_at"`
	OrderItems []OrderItemEntity    `json:"order_items" database:"order_items"`
	Payments   []OrderPaymentEntity `json:"payments" database:"payments"`
}
