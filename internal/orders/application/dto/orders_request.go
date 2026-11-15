package dto

import "microservice/internal/orders/domain/entity"

type OrderUpdateRequest struct {
	Id         string                   `json:"id" example:"order-123"`
	Rut        string                   `json:"rut" example:"12345678-9"`
	SubTotal   float64                  `json:"sub_total" example:"12500"`
	OrderItems []OrderItemUpdateRequest `json:"order_items"`
}

type OrderItemUpdateRequest struct {
	Sku      string `json:"sku" example:"SKU-1"`
	Quantity int    `json:"quantity" example:"2"`
}

func NewOrderUpdateRequest() *OrderUpdateRequest {
	response := &OrderUpdateRequest{
		OrderItems: []OrderItemUpdateRequest{},
	}
	return response
}

func (orderUpdateRequest *OrderUpdateRequest) ToOrderEntity() (*entity.OrderEntity, error) {
	orderEntity := &entity.OrderEntity{
		Id:         orderUpdateRequest.Id,
		Rut:        orderUpdateRequest.Rut,
		SubTotal:   orderUpdateRequest.SubTotal,
		OrderItems: []entity.OrderItemEntity{},
	}

	for _, item := range orderUpdateRequest.OrderItems {
		orderItem := entity.OrderItemEntity{
			OrderId:  orderUpdateRequest.Id,
			Sku:      item.Sku,
			Quantity: item.Quantity,
		}
		orderEntity.OrderItems = append(orderEntity.OrderItems, orderItem)
	}

	return orderEntity, nil
}
