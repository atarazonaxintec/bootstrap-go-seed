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

// @Summary Create Order
// @Description Creates a new order for the given RUT
// @Tags Orders
// @Accept json
// @Produce json
// @Param rut path string true "RUT of the customer"
// @Success 201 {object} dto.OrdersResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /orders/{rut} [post]
func (h *OrdersHandler) CreateOrder(c *gin.Context, rut string) {
	res, err := h.uc.NewOrders(c.Request.Context(), rut)

	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary Update Order
// @Description Updates an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dto.OrderUpdateRequest true "Order data to update"
// @Success 200 {object} dto.OrdersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /orders [put]
func (h *OrdersHandler) UpdateOrder(c *gin.Context, order *dto.OrderUpdateRequest) {
	orderEntity, err := order.ToOrderEntity()
	if err != nil {
		errorResponse := dto.NewErrorResponse("Invalid order data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	res, err := h.uc.UpdateOrders(c.Request.Context(), orderEntity)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Get Order by ID
// @Description Retrieves an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrdersResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrdersHandler) GetOrderById(c *gin.Context, id string) {
	res, err := h.uc.GetOrderById(c.Request.Context(), id)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	if res == nil {
		errorResponse := dto.NewErrorResponse("Order not found")
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Get Orders by RUT
// @Description Retrieves all orders for a given RUT
// @Tags Orders
// @Accept json
// @Produce json
// @Param rut path string true "RUT of the customer"
// @Success 200 {array} dto.OrdersResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /orders/rut/{rut} [get]
func (h *OrdersHandler) GetOrdersByRut(c *gin.Context, rut string) {
	res, err := h.uc.GetOrdersByRut(c.Request.Context(), rut)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	if res == nil {
		errorResponse := dto.NewErrorResponse("No orders found for the given RUT")
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
