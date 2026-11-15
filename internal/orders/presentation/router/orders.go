package router

import (
	dto "microservice/internal/orders/application/dto"
	"microservice/internal/orders/presentation/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.OrdersHandler) {
	orders := api.Group("/orders")
	{
		orders.POST("/:rut", func(c *gin.Context) {
			rut := c.Param("rut")
			h.CreateOrder(c, rut)
		})

		orders.PUT("", func(c *gin.Context) {
			var orderUpdateRequest dto.OrderUpdateRequest
			if err := c.ShouldBindJSON(&orderUpdateRequest); err != nil {
				c.JSON(400, gin.H{"error": "Invalid request body"})
				return
			}
			h.UpdateOrder(c, &orderUpdateRequest)
		})

		orders.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			h.GetOrderById(c, id)
		})

		orders.GET("/rut/:rut", func(c *gin.Context) {
			rut := c.Param("rut")
			h.GetOrdersByRut(c, rut)
		})
	}
}
