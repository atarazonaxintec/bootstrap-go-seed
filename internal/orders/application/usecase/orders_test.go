package usecase

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/orders/domain/entity"
)

type mockOrdersRepository struct {
	createdOrder        *entity.OrderEntity
	createErr           error
	receivedRut         string
	updatedOrder        *entity.OrderEntity
	updateErr           error
	receivedOrder       *entity.OrderEntity
	orderByID           *entity.OrderEntity
	getOrderByIDErr     error
	receivedID          string
	ordersByRut         []*entity.OrderEntity
	getOrdersByRutErr   error
	receivedOrdersByRut string
}

func (m *mockOrdersRepository) Ping(context.Context) error {
	return nil
}

func (m *mockOrdersRepository) CreateOrder(_ context.Context, rut string) (*entity.OrderEntity, error) {
	m.receivedRut = rut
	return m.createdOrder, m.createErr
}

func (m *mockOrdersRepository) UpdateOrder(_ context.Context, order *entity.OrderEntity) (*entity.OrderEntity, error) {
	m.receivedOrder = order
	return m.updatedOrder, m.updateErr
}

func (m *mockOrdersRepository) GetOrderById(_ context.Context, id string) (*entity.OrderEntity, error) {
	m.receivedID = id
	return m.orderByID, m.getOrderByIDErr
}

func (m *mockOrdersRepository) GetOrdersByRut(_ context.Context, rut string) ([]*entity.OrderEntity, error) {
	m.receivedOrdersByRut = rut
	return m.ordersByRut, m.getOrdersByRutErr
}

func (m *mockOrdersRepository) GetOrderAndReserveForUpdate(context.Context, string) (*entity.OrderEntity, error) {
	return nil, nil
}

func TestOrdersUseCase_NewOrders(t *testing.T) {
	createdOrder := &entity.OrderEntity{
		Id:       "order-123",
		Rut:      "12345678-9",
		SubTotal: 12500,
		Status:   entity.OrderStatus("en_preparacion"),
		OrderItems: []entity.OrderItemEntity{
			{OrderId: "order-123", Sku: "SKU-1", Quantity: 2},
		},
		Payments: []entity.OrderPaymentEntity{
			{
				Id:            "payment-123",
				PaymentMethod: "card",
				Amount:        12500,
				Status:        entity.OrderPaymentStatus("approved"),
			},
		},
	}
	mockRepo := &mockOrdersRepository{createdOrder: createdOrder}
	uc := NewOrdersUseCase(mockRepo)

	res, err := uc.NewOrders(context.Background(), createdOrder.Rut)
	if err != nil {
		t.Fatalf("NewOrders() returned unexpected error: %v", err)
	}
	if mockRepo.receivedRut != createdOrder.Rut {
		t.Errorf("CreateOrder() received RUT %q, want %q", mockRepo.receivedRut, createdOrder.Rut)
	}
	if res == nil {
		t.Fatal("NewOrders() returned nil response")
	}
	if res.Id != createdOrder.Id || res.Rut != createdOrder.Rut || res.Status != createdOrder.Status || res.SubTotal != createdOrder.SubTotal {
		t.Errorf("NewOrders() response does not match created order: %+v", res)
	}
	if len(res.OrderResponseItems) != 1 || res.OrderResponseItems[0].Sku != "SKU-1" || res.OrderResponseItems[0].Quantity != 2 {
		t.Errorf("NewOrders() order items = %+v, want the created order item", res.OrderResponseItems)
	}
	if len(res.PaymentResponseItems) != 1 || res.PaymentResponseItems[0].Id != "payment-123" || res.PaymentResponseItems[0].Amount != 12500 {
		t.Errorf("NewOrders() payments = %+v, want the created order payment", res.PaymentResponseItems)
	}
}

func TestOrdersUseCase_NewOrders_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("create order failed")
	uc := NewOrdersUseCase(&mockOrdersRepository{createErr: repositoryErr})

	res, err := uc.NewOrders(context.Background(), "12345678-9")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("NewOrders() error = %v, want %v", err, repositoryErr)
	}
	if res != nil {
		t.Errorf("NewOrders() response = %+v, want nil on error", res)
	}
}

func TestOrdersUseCase_UpdateOrders(t *testing.T) {
	inputOrder := &entity.OrderEntity{Id: "order-123", Rut: "12345678-9"}
	updatedOrder := &entity.OrderEntity{
		Id:       inputOrder.Id,
		Rut:      inputOrder.Rut,
		SubTotal: 12500,
		Status:   entity.OrderStatus("en_transito"),
	}
	mockRepo := &mockOrdersRepository{updatedOrder: updatedOrder}
	uc := NewOrdersUseCase(mockRepo)

	res, err := uc.UpdateOrders(context.Background(), inputOrder)
	if err != nil {
		t.Fatalf("UpdateOrders() returned unexpected error: %v", err)
	}
	if mockRepo.receivedOrder != inputOrder {
		t.Fatal("UpdateOrder() did not receive the input order")
	}
	if res == nil {
		t.Fatal("UpdateOrders() returned nil response")
	}
	if res.Id != updatedOrder.Id || res.Rut != updatedOrder.Rut || res.Status != updatedOrder.Status || res.SubTotal != updatedOrder.SubTotal {
		t.Errorf("UpdateOrders() response does not match updated order: %+v", res)
	}
}

func TestOrdersUseCase_UpdateOrders_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("update order failed")
	uc := NewOrdersUseCase(&mockOrdersRepository{updateErr: repositoryErr})

	res, err := uc.UpdateOrders(context.Background(), &entity.OrderEntity{Id: "order-123"})
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("UpdateOrders() error = %v, want %v", err, repositoryErr)
	}
	if res != nil {
		t.Errorf("UpdateOrders() response = %+v, want nil on error", res)
	}
}

func TestOrdersUseCase_GetOrderById(t *testing.T) {
	order := &entity.OrderEntity{
		Id:       "order-123",
		Rut:      "12345678-9",
		SubTotal: 12500,
		Status:   entity.OrderStatus("en_transito"),
	}
	mockRepo := &mockOrdersRepository{orderByID: order}
	uc := NewOrdersUseCase(mockRepo)

	res, err := uc.GetOrderById(context.Background(), order.Id)
	if err != nil {
		t.Fatalf("GetOrderById() returned unexpected error: %v", err)
	}
	if mockRepo.receivedID != order.Id {
		t.Errorf("GetOrderById() queried ID %q, want %q", mockRepo.receivedID, order.Id)
	}
	if res == nil {
		t.Fatal("GetOrderById() returned nil response")
	}
	if res.Id != order.Id || res.Rut != order.Rut || res.Status != order.Status || res.SubTotal != order.SubTotal {
		t.Errorf("GetOrderById() response does not match order: %+v", res)
	}
}

func TestOrdersUseCase_GetOrderById_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("get order failed")
	uc := NewOrdersUseCase(&mockOrdersRepository{getOrderByIDErr: repositoryErr})

	res, err := uc.GetOrderById(context.Background(), "order-123")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("GetOrderById() error = %v, want %v", err, repositoryErr)
	}
	if res != nil {
		t.Errorf("GetOrderById() response = %+v, want nil on error", res)
	}
}

func TestOrdersUseCase_GetOrdersByRut(t *testing.T) {
	rut := "12345678-9"
	orders := []*entity.OrderEntity{
		{Id: "order-1", Rut: rut, SubTotal: 1000},
		{Id: "order-2", Rut: rut, SubTotal: 2000},
	}
	mockRepo := &mockOrdersRepository{ordersByRut: orders}
	uc := NewOrdersUseCase(mockRepo)

	res, err := uc.GetOrdersByRut(context.Background(), rut)
	if err != nil {
		t.Fatalf("GetOrdersByRut() returned unexpected error: %v", err)
	}
	if mockRepo.receivedOrdersByRut != rut {
		t.Errorf("GetOrdersByRut() queried RUT %q, want %q", mockRepo.receivedOrdersByRut, rut)
	}
	if len(res) != len(orders) {
		t.Fatalf("GetOrdersByRut() returned %d orders, want %d", len(res), len(orders))
	}
	for i, order := range orders {
		if res[i].Id != order.Id || res[i].Rut != order.Rut || res[i].SubTotal != order.SubTotal {
			t.Errorf("GetOrdersByRut() response[%d] = %+v, want order %q", i, res[i], order.Id)
		}
	}
}

func TestOrdersUseCase_GetOrdersByRut_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("get orders failed")
	uc := NewOrdersUseCase(&mockOrdersRepository{getOrdersByRutErr: repositoryErr})

	res, err := uc.GetOrdersByRut(context.Background(), "12345678-9")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("GetOrdersByRut() error = %v, want %v", err, repositoryErr)
	}
	if res != nil {
		t.Errorf("GetOrdersByRut() response = %+v, want nil on error", res)
	}
}
