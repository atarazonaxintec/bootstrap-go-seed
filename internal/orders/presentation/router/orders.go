package router

import (
	"microservice/internal/orders/presentation/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.OrdersHandler) {
	v1 := api.Group("/v1/orders")
	v2 := api.Group("/v2/orders")

	v1.GET("", h.Orders)
	v2.GET("", h.Orders)
}
