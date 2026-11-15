package handler

import (
	"microservice/internal/orders/application/dto"
	"microservice/internal/orders/application/usecase"
	domain "microservice/internal/orders/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentsHandler struct {
	uc *usecase.PaymentsUseCase
}

// Constructor del Controller
func NewPaymentsHandler(uc *usecase.PaymentsUseCase) *PaymentsHandler {
	return &PaymentsHandler{
		uc: uc,
	}
}

// @Summary Create Empty Payment
// @Description Creates an empty payment for the given order ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 201 {object} entity.OrderPaymentEntity
// @Failure 503 {object} dto.ErrorResponse
// @Router /payments/{orderID} [post]
func (h *PaymentsHandler) CreatePayment(c *gin.Context, orderID string) {
	payment, err := h.uc.CreateEmptyPayment(c.Request.Context(), orderID)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// @Summary Get Payments by Order ID
// @Description Retrieves a list of payments for the given order ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} []dto.OrderPaymentResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /payments/order/{orderID} [get]
func (h *PaymentsHandler) GetPaymentsByOrderID(c *gin.Context, orderID string) {
	payments, err := h.uc.FindPaymentsByOrderID(c.Request.Context(), orderID)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}
	if len(payments) == 0 {
		errorResponse := dto.NewErrorResponse("No payments found for the given order ID")
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, payments)
}

// @Summary Pay for an Order
// @Description Processes a payment for the given order ID, payment method, and amount
// @Tags Payments
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Param paymentMethod path string true "Payment Method"
// @Param amount path number true "Payment Amount"
// @Success 200 {object} entity.OrderPaymentEntity
// @Failure 503 {object} dto.ErrorResponse
// @Router /payments/pay/{orderID}/{paymentMethod}/{amount} [post]
func (h *PaymentsHandler) Pay(c *gin.Context, orderID string, paymentMethod string, amount float64) {
	payment, err := h.uc.Pay(c.Request.Context(), orderID, paymentMethod, amount)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentsHandler) UpdatePaymentStatus(c *gin.Context, orderID string, status domain.OrderPaymentStatus) {
	err := h.uc.UpdatePaymentStatus(c.Request.Context(), orderID, status)
	if err != nil {
		errorResponse := dto.NewErrorResponse("Service unavailable")
		c.JSON(http.StatusServiceUnavailable, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}
