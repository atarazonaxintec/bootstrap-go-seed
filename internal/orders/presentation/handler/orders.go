package handler

import (
	"microservice/internal/orders/application/dto"
	"microservice/internal/orders/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	uc *usecase.OrdersUseCase
}

// Constructor del Controller
func NewOrdersHandler(uc *usecase.OrdersUseCase) *OrdersHandler {
	return &OrdersHandler{
		uc: uc,
	}
}

// @Summary Orders Check
// @Description Returns the orders status of the service
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {object} dto.OrdersResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /orders [get]
func (h *OrdersHandler) Orders(c *gin.Context) {
	res, err := h.uc.Execute(c.Request.Context())

	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unordersy")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
