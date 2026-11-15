package usecase

import (
	"context"
	"errors"
	"microservice/internal/orders/domain/entity"
	"microservice/pkg/logger"
	"os"
	"testing"
)

type mockProductsRepository struct {
	products     []*entity.ProductEntity
	findAllErr   error
	product      *entity.ProductEntity
	findBySkuErr error
	receivedSku  string
}

func TestMain(m *testing.M) {
	logger.Init("dev")
	os.Exit(m.Run())
}

func (m *mockProductsRepository) FindAll(context.Context) ([]*entity.ProductEntity, error) {
	return m.products, m.findAllErr
}

func (m *mockProductsRepository) FindBySku(_ context.Context, sku string) (*entity.ProductEntity, error) {
	m.receivedSku = sku
	return m.product, m.findBySkuErr
}

func TestProductsUseCase_FindAll(t *testing.T) {
	repo := &mockProductsRepository{
		products: []*entity.ProductEntity{
			{Sku: "SKU-1", Name: "Product 1", Description: "First", Price: 1250, StockQuantity: 3, CreatedAt: "created-1", UpdatedAt: "updated-1"},
			{Sku: "SKU-2", Name: "Product 2", Description: "Second", Price: 2500, StockQuantity: 5, CreatedAt: "created-2", UpdatedAt: "updated-2"},
		},
	}
	uc := NewProductsUseCase(repo)

	got, err := uc.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() returned unexpected error: %v", err)
	}
	if len(got) != len(repo.products) {
		t.Fatalf("FindAll() returned %d products, want %d", len(got), len(repo.products))
	}
	for i, product := range repo.products {
		response := got[i]
		if response.Sku != product.Sku ||
			response.Name != product.Name ||
			response.Description != product.Description ||
			response.Price != product.Price ||
			response.StockQuantity != product.StockQuantity ||
			response.CreatedAt != product.CreatedAt ||
			response.UpdatedAt != product.UpdatedAt {
			t.Errorf("FindAll() response[%d] = %+v, want product %+v", i, response, product)
		}
	}
}

func TestProductsUseCase_FindAll_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("find products failed")
	uc := NewProductsUseCase(&mockProductsRepository{findAllErr: repositoryErr})

	got, err := uc.FindAll(context.Background())
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("FindAll() error = %v, want %v", err, repositoryErr)
	}
	if got != nil {
		t.Errorf("FindAll() = %+v, want nil on error", got)
	}
}

func TestProductsUseCase_FindBySku(t *testing.T) {
	product := &entity.ProductEntity{
		Sku:           "SKU-1",
		Name:          "Product 1",
		Description:   "First",
		Price:         1250,
		StockQuantity: 3,
		CreatedAt:     "created",
		UpdatedAt:     "updated",
	}
	repo := &mockProductsRepository{product: product}
	uc := NewProductsUseCase(repo)

	got, err := uc.FindBySku(context.Background(), product.Sku)
	if err != nil {
		t.Fatalf("FindBySku() returned unexpected error: %v", err)
	}
	if repo.receivedSku != product.Sku {
		t.Errorf("FindBySku() queried SKU %q, want %q", repo.receivedSku, product.Sku)
	}
	if got == nil {
		t.Fatal("FindBySku() returned nil response")
	}
	if got.Sku != product.Sku ||
		got.Name != product.Name ||
		got.Description != product.Description ||
		got.Price != product.Price ||
		got.StockQuantity != product.StockQuantity ||
		got.CreatedAt != product.CreatedAt ||
		got.UpdatedAt != product.UpdatedAt {
		t.Errorf("FindBySku() response = %+v, want product %+v", got, product)
	}
}

func TestProductsUseCase_FindBySku_ReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("find product failed")
	uc := NewProductsUseCase(&mockProductsRepository{findBySkuErr: repositoryErr})

	got, err := uc.FindBySku(context.Background(), "SKU-1")
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("FindBySku() error = %v, want %v", err, repositoryErr)
	}
	if got != nil {
		t.Errorf("FindBySku() = %+v, want nil on error", got)
	}
}
