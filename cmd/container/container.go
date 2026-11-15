package container

import (
	"microservice/cmd/config"
	healthUseCase "microservice/internal/health/application/usecase"
	healthInfra "microservice/internal/health/infrastructure"
	healthHandler "microservice/internal/health/presentation/handler"
	ordersUseCase "microservice/internal/orders/application/usecase"
	ordersInfra "microservice/internal/orders/infrastructure"
	ordersHandler "microservice/internal/orders/presentation/handler"
	"microservice/pkg/database"
	"microservice/pkg/logger"
	"microservice/pkg/pubsub"

	"go.uber.org/zap"
)

type Container struct {
	DB              database.Database
	PubSub          pubsub.PubSub
	HealthHandler   *healthHandler.HealthHandler
	OrdersHandler   *ordersHandler.OrdersHandler
	PaymentsHandler *ordersHandler.PaymentsHandler
	ProductsHandler *ordersHandler.ProductsHandler
	Log             *zap.Logger
}

func NewContainer(cfg *config.Configuration, db database.Database, ps pubsub.PubSub, log *zap.Logger) (*Container, error) {
	healthRepo := healthInfra.NewPostgresHealthRepository(db)
	healthUC := healthUseCase.NewHealthUseCase(healthRepo)
	healthHdlr := healthHandler.NewHealthHandler(healthUC)

	ordersRepo := ordersInfra.NewPostgresOrdersRepository(db)
	ordersUC := ordersUseCase.NewOrdersUseCase(ordersRepo)
	ordersHdlr := ordersHandler.NewOrdersHandler(ordersUC)

	productsRepo := ordersInfra.NewPostgresProductsRepository(db)
	productsUC := ordersUseCase.NewProductsUseCase(productsRepo)
	productsHdlr := ordersHandler.NewProductsHandler(productsUC)

	paymentsRepo := ordersInfra.NewPostgresPaymentsRepository(db)
	paymentsUC := ordersUseCase.NewPaymentsUseCase(paymentsRepo)
	paymentsHdlr := ordersHandler.NewPaymentsHandler(paymentsUC)

	logger.Success("Container de dependencias inicializado...")

	return &Container{
		DB:              db,
		PubSub:          ps,
		HealthHandler:   healthHdlr,
		OrdersHandler:   ordersHdlr,
		PaymentsHandler: paymentsHdlr,
		ProductsHandler: productsHdlr,
		Log:             log,
	}, nil
}
