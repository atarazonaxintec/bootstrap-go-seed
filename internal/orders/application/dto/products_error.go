package dto

import "time"

// ProductsErrorResponse represents a products error response.
type ProductsErrorResponse struct {
	Error     string    `json:"error" example:"Service unordersy"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

// NewProductsErrorResponse creates a products error response.
func NewProductsErrorResponse(message string) *ProductsErrorResponse {
	return &ProductsErrorResponse{
		Error:     message,
		Timestamp: time.Now(),
	}
}
