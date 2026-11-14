package dto

import (
	entity "microservice/internal/orders/domain/entity"
)

type OrdersResponse struct {
	Id                   string                `json:"id"`
	Status               entity.OrderStatus    `json:"status"`
	Rut                  string                `json:"rut"`
	SubTotal             float64               `json:"sub_total"`
	OrderResponseItems   []OrderResponseItem   `json:"order_items"`
	PaymentResponseItems []PaymentResponseItem `json:"payments"`
}

type OrderResponseItem struct {
	Id       string  `json:"id"`
	Sku      string  `json:"product_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type PaymentResponseItem struct {
	Id     string                    `json:"id"`
	Amount float64                   `json:"amount"`
	Method string                    `json:"method"`
	Status entity.OrderPaymentStatus `json:"status"`
}

// Constructor
func NewOrdersResponse() *OrdersResponse {
	response := &OrdersResponse{
		OrderResponseItems:   []OrderResponseItem{},
		PaymentResponseItems: []PaymentResponseItem{},
	}
	return response
}

func EditOrderResponse(entityOrder *entity.Order) *OrdersResponse {
	response := NewOrdersResponse()
	response.Id = entityOrder.Id
	response.Status = entityOrder.Status
	response.Rut = entityOrder.Rut
	response.SubTotal = entityOrder.SubTotal

	for _, item := range entityOrder.OrderItems {
		orderItem := OrderResponseItem{
			Id:       item.OrderId,
			Sku:      item.Sku,
			Quantity: item.Quantity,
		}
		response.OrderResponseItems = append(response.OrderResponseItems, orderItem)
	}

	for _, payment := range entityOrder.Payments {
		paymentItem := PaymentResponseItem{
			Id:     payment.Id,
			Amount: payment.Amount,
			Method: payment.PaymentMethod,
			Status: payment.Status,
		}
		response.PaymentResponseItems = append(response.PaymentResponseItems, paymentItem)
	}

	return response
}
