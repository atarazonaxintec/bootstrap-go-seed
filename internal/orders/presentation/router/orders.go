package router

import (
	dto "microservice/internal/orders/application/dto"
	"microservice/internal/orders/presentation/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, hOrders *handler.OrdersHandler, hPayments *handler.PaymentsHandler, hProducts *handler.ProductsHandler) {
	orders := api.Group("/orders")
	{
		orders.POST("/rut/:rut", func(c *gin.Context) {
			rut := c.Param("rut")
			hOrders.CreateOrder(c, rut)
		})

		orders.PUT("", func(c *gin.Context) {
			var orderUpdateRequest dto.OrderUpdateRequest
			if err := c.ShouldBindJSON(&orderUpdateRequest); err != nil {
				c.JSON(400, gin.H{"error": "Invalid request body"})
				return
			}
			hOrders.UpdateOrder(c, &orderUpdateRequest)
		})

		orders.GET("/id/:id", func(c *gin.Context) {
			id := c.Param("id")
			hOrders.GetOrderById(c, id)
		})

		orders.GET("/rut/:rut", func(c *gin.Context) {
			rut := c.Param("rut")
			hOrders.GetOrdersByRut(c, rut)
		})
	}

	payments := api.Group("/payments")
	{
		payments.POST("/:orderID", func(c *gin.Context) {
			orderID := c.Param("orderID")
			hPayments.CreatePayment(c, orderID)
		})

		payments.GET("/order/:orderID", func(c *gin.Context) {
			orderID := c.Param("orderID")
			hPayments.GetPaymentsByOrderID(c, orderID)
		})
	}

	products := api.Group("/products")
	{
		products.GET("", func(c *gin.Context) {
			hProducts.GetAllProducts(c)
		})

		products.GET("/:sku", func(c *gin.Context) {
			sku := c.Param("sku")
			hProducts.GetProductBySku(c, sku)
		})
	}
}
