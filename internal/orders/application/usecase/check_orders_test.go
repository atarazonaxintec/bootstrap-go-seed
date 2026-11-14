package usecase

import (
	"context"
	"testing"
)

type mockOrdersRepository struct{}

func (m *mockOrdersRepository) Ping(ctx context.Context) error {
	return nil
}

func TestOrdersUseCase(t *testing.T) {
	mockRepo := &mockOrdersRepository{}
	// Patrón
	// 1) Arrancamos creando el use case
	uc := NewOrdersUseCase(mockRepo)

	// 2) Ejecutamos el metodo
	res, err := uc.Execute(context.Background())

	// 3) Validamos resultados
	if err != nil {
		t.Fatalf("Execute() retuned unexpected error: %v", err)
	}

	if res == nil {
		t.Fatalf("Exceute retuned nil response")
	}

	if res.Service != "microservice" {
		t.Errorf("Execute().Service = %q, want %q", res.Service, "microservice")
	}

	if res.Timestamp.IsZero() {
		t.Errorf("Execute().Timestamp is zero, expected a valid timestamp")
	}

	if res.Database != "connected" {
		t.Errorf("Execute().Database = %q, want %q", res.Database, "connected")
	}
}

// Ejemplo cuando el use case tenga un repositorio:
// func TestOrdersUseCase_Execute_WithDB(t *testing.T) {
//     mockRepo := new(MockOrdersRepository)
//     mockRepo.On("Check").Return(nil)
//     uc := NewOrdersUseCase(mockRepo)
//     res, err := uc.Execute()
//     assert.NoError(t, err)
//     assert.NotNil(t, res)
// }
