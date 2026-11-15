package handler

import (
	"microservice/internal/orders/application/dto"
	"microservice/internal/orders/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductsHandler struct {
	uc *usecase.ProductsUseCase
}

// Constructor del Controller
func NewProductsHandler(uc *usecase.ProductsUseCase) *ProductsHandler {
	return &ProductsHandler{
		uc: uc,
	}
}

// @Summary Get All Products
// @Description Retrieves a list of all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ProductsResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /products [get]
func (h *ProductsHandler) GetAllProducts(c *gin.Context) {
	res, err := h.uc.FindAll(c.Request.Context())
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	if len(res) == 0 {
		errorResponse := dto.NewErrorResponse("No products found")
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Get Product by SKU
// @Description Retrieves a product by its SKU
// @Tags Products
// @Accept json
// @Produce json
// @Param sku path string true "Product SKU"
// @Success 200 {object} dto.ProductsResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /products/{sku} [get]
func (h *ProductsHandler) GetProductBySku(c *gin.Context, sku string) {
	res, err := h.uc.FindBySku(c.Request.Context(), sku)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	if res == nil {
		errorResponse := dto.NewErrorResponse("Product not found")
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
