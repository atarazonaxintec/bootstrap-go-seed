package usecase

import (
	"context"
	"errors"
	"testing"

	"microservice/internal/orders/domain/entity"
)

type mockPaymentRepository struct {
	emptyPayment        *entity.OrderPaymentEntity
	createEmptyErr      error
	createdOrderID      string
	payments            []entity.OrderPaymentEntity
	findPaymentsErr     error
	foundOrderID        string
	paidPayment         *entity.OrderPaymentEntity
	payErr              error
	paidOrderID         string
	paidMethod          string
	paidAmount          float64
	updateStatusErr     error
	updatedOrderID      string
	updatedPaymentState entity.OrderPaymentStatus
}

func (m *mockPaymentRepository) CreateEmpty(_ context.Context, orderID string) (*entity.OrderPaymentEntity, error) {
	m.createdOrderID = orderID
	return m.emptyPayment, m.createEmptyErr
}

func (m *mockPaymentRepository) FindByOrderID(_ context.Context, orderID string) ([]entity.OrderPaymentEntity, error) {
	m.foundOrderID = orderID
	return m.payments, m.findPaymentsErr
}

func (m *mockPaymentRepository) Pay(_ context.Context, orderID, paymentMethod string, amount float64) (*entity.OrderPaymentEntity, error) {
	m.paidOrderID = orderID
	m.paidMethod = paymentMethod
	m.paidAmount = amount
	return m.paidPayment, m.payErr
}

func (m *mockPaymentRepository) UpdateStatus(_ context.Context, orderID string, status entity.OrderPaymentStatus) error {
	m.updatedOrderID = orderID
	m.updatedPaymentState = status
	return m.updateStatusErr
}

func TestPaymentsUseCase_CreateEmptyPayment(t *testing.T) {
	payment := &entity.OrderPaymentEntity{Id: "payment-1", OrderId: "order-1"}
	repo := &mockPaymentRepository{emptyPayment: payment}
	uc := NewPaymentsUseCase(repo)

	got, err := uc.CreateEmptyPayment(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("CreateEmptyPayment() returned unexpected error: %v", err)
	}
	if repo.createdOrderID != "order-1" {
		t.Errorf("CreateEmpty() order ID = %q, want %q", repo.createdOrderID, "order-1")
	}
	if got != payment {
		t.Errorf("CreateEmptyPayment() = %p, want %p", got, payment)
	}
}

func TestPaymentsUseCase_CreateEmptyPayment_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("create empty payment failed")
	uc := NewPaymentsUseCase(&mockPaymentRepository{createEmptyErr: repositoryErr})

	got, err := uc.CreateEmptyPayment(context.Background(), "order-1")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("CreateEmptyPayment() error = %v, want %v", err, repositoryErr)
	}
	if got != nil {
		t.Errorf("CreateEmptyPayment() = %+v, want nil on error", got)
	}
}

func TestPaymentsUseCase_FindPaymentsByOrderID(t *testing.T) {
	repo := &mockPaymentRepository{
		payments: []entity.OrderPaymentEntity{
			{
				Id:            "payment-1",
				PaymentMethod: "card",
				Status:        entity.OrderPaymentStatus("approved"),
				Amount:        1250,
				PaidAt:        "2026-11-15T10:00:00Z",
			},
		},
	}
	uc := NewPaymentsUseCase(repo)

	got, err := uc.FindPaymentsByOrderID(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("FindPaymentsByOrderID() returned unexpected error: %v", err)
	}
	if repo.foundOrderID != "order-1" {
		t.Errorf("FindByOrderID() order ID = %q, want %q", repo.foundOrderID, "order-1")
	}
	if len(got) != 1 {
		t.Fatalf("FindPaymentsByOrderID() returned %d payments, want 1", len(got))
	}
	payment := got[0]
	if payment.PaymentID != "payment-1" || payment.PaymentMethod != "card" || payment.Status != "approved" || payment.Amount != 1250 || payment.PaidAt != "2026-11-15T10:00:00Z" {
		t.Errorf("FindPaymentsByOrderID() returned unexpected payment: %+v", payment)
	}
}

func TestPaymentsUseCase_FindPaymentsByOrderID_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("find payments failed")
	uc := NewPaymentsUseCase(&mockPaymentRepository{findPaymentsErr: repositoryErr})

	got, err := uc.FindPaymentsByOrderID(context.Background(), "order-1")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("FindPaymentsByOrderID() error = %v, want %v", err, repositoryErr)
	}
	if got != nil {
		t.Errorf("FindPaymentsByOrderID() = %+v, want nil on error", got)
	}
}

func TestPaymentsUseCase_Pay(t *testing.T) {
	payment := &entity.OrderPaymentEntity{Id: "payment-1", OrderId: "order-1", Amount: 1250}
	repo := &mockPaymentRepository{paidPayment: payment}
	uc := NewPaymentsUseCase(repo)

	got, err := uc.Pay(context.Background(), "order-1", "card", 1250)
	if err != nil {
		t.Fatalf("Pay() returned unexpected error: %v", err)
	}
	if repo.paidOrderID != "order-1" || repo.paidMethod != "card" || repo.paidAmount != 1250 {
		t.Errorf("Pay() arguments = (%q, %q, %v), want (%q, %q, %v)", repo.paidOrderID, repo.paidMethod, repo.paidAmount, "order-1", "card", 1250.0)
	}
	if got != payment {
		t.Errorf("Pay() = %p, want %p", got, payment)
	}
}

func TestPaymentsUseCase_Pay_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("pay failed")
	uc := NewPaymentsUseCase(&mockPaymentRepository{payErr: repositoryErr})

	got, err := uc.Pay(context.Background(), "order-1", "card", 1250)
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("Pay() error = %v, want %v", err, repositoryErr)
	}
	if got != nil {
		t.Errorf("Pay() = %+v, want nil on error", got)
	}
}

func TestPaymentsUseCase_UpdatePaymentStatus(t *testing.T) {
	repo := &mockPaymentRepository{}
	uc := NewPaymentsUseCase(repo)
	status := entity.OrderPaymentStatus("approved")

	if err := uc.UpdatePaymentStatus(context.Background(), "order-1", status); err != nil {
		t.Fatalf("UpdatePaymentStatus() returned unexpected error: %v", err)
	}
	if repo.updatedOrderID != "order-1" || repo.updatedPaymentState != status {
		t.Errorf("UpdateStatus() arguments = (%q, %q), want (%q, %q)", repo.updatedOrderID, repo.updatedPaymentState, "order-1", status)
	}
}

func TestPaymentsUseCase_UpdatePaymentStatus_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("update payment status failed")
	uc := NewPaymentsUseCase(&mockPaymentRepository{updateStatusErr: repositoryErr})

	err := uc.UpdatePaymentStatus(context.Background(), "order-1", entity.OrderPaymentStatus("approved"))
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("UpdatePaymentStatus() error = %v, want %v", err, repositoryErr)
	}
}
