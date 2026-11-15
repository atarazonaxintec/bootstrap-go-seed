package routes

import (
	"microservice/cmd/container"

	health "microservice/internal/health/presentation/router"
	orders "microservice/internal/orders/presentation/router"

	"github.com/gin-gonic/gin"
)

func ConfigureModule(router *gin.Engine, c *container.Container) {
	api := router.Group("api")
	health.RegisterRoutes(api, c.HealthHandler)
	// Registramos las rutas de forma dinámica con su versión desde aquí
	// Orders
	orders.RegisterRoutes(api, c.OrdersHandler, c.PaymentsHandler, c.ProductsHandler)
}
