package router

import (
	"github.com/gofiber/fiber/v2"

	order_handler "github.com/developer-afo/instashop-ecommerce-api/handler/order"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/middleware"
	coreRepository "github.com/developer-afo/instashop-ecommerce-api/repository/core"
	finance_repository "github.com/developer-afo/instashop-ecommerce-api/repository/finance"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
	user_repository "github.com/developer-afo/instashop-ecommerce-api/repository/user"
	"github.com/developer-afo/instashop-ecommerce-api/service"
	core_service "github.com/developer-afo/instashop-ecommerce-api/service/core"
	finance_service "github.com/developer-afo/instashop-ecommerce-api/service/finance"
	payment_gateway_service "github.com/developer-afo/instashop-ecommerce-api/service/finance/payment_gateway"
	order_service "github.com/developer-afo/instashop-ecommerce-api/service/order"
	user_service "github.com/developer-afo/instashop-ecommerce-api/service/user"
)

func InitializeOrderRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)
	orderRepository := order_repository.NewOrderRepository(db)
	orderItemRepository := order_repository.NewOrderItemRepository(db)
	orderStatusRepository := order_repository.NewOrderStatusRepository(db)
	orderStatusHistoryRepository := order_repository.NewOrderStatusHistoryRepository(db)
	imageRepository := coreRepository.NewImageRepository(db)
	productRepository := coreRepository.NewProductRepository(db)
	transactionRepository := finance_repository.NewTransactionRepository(db)

	// Services
	httpService := service.NewHTTPService()

	imageService := core_service.NewImageService(imageRepository)
	productService := core_service.NewProductService(productRepository, imageService)

	orderItemService := order_service.NewOrderItemService(orderItemRepository, productService)
	orderStatusService := order_service.NewOrderStatusService(orderStatusRepository)
	orderStatusHistoryService := order_service.NewOrderStatusHistoryService(orderStatusHistoryRepository, orderStatusService)

	transactionService := finance_service.NewTransactionService(transactionRepository)
	paystackPaymentService := payment_gateway_service.NewPaystackService(httpService, env)
	flutterwavePaymentService := payment_gateway_service.NewFlutterwaveService(httpService, env)
	paymentGatewayService := payment_gateway_service.NewPaymentGatewayService(paystackPaymentService, flutterwavePaymentService)

	userService := user_service.NewUserService(userRepository)

	orderService := order_service.NewOrderService(
		orderRepository,
		orderItemService,
		orderStatusService,
		orderStatusHistoryService,
		productService,
		transactionService,
		paymentGatewayService,
		userService,
	)

	// Handlers
	orderHandler := order_handler.NewOrderHandler(orderService, orderStatusHistoryService, orderStatusService)

	// middlewares
	roleMiddleware := middleware.NewRoleMiddleware(userRepository)
	authMiddleware := middleware.Protected()

	// Base routes
	orderRouter := router.Group("/order", authMiddleware)

	// Routes
	orderRouter.Post("/", roleMiddleware.ValidateRole(user_service.UserRoleCustomer), orderHandler.CreateOrder)
	orderRouter.Post("/cancel/:order_id", roleMiddleware.ValidateRole(user_service.UserRoleCustomer), orderHandler.CancelOrder)
	orderRouter.Get("/", roleMiddleware.ValidateRole(user_service.UserRoleCustomer), orderHandler.GetUserOrders)
	orderRouter.Get("/all", roleMiddleware.ValidateRole(user_service.UserRoleAdmin), orderHandler.GetAllOrders)
	orderRouter.Get("/verify-payment/:reference", orderHandler.VerifyOrderPayment)
	orderRouter.Get("/status-history/:order_id", orderHandler.StatusHistoryByOrderId)
	orderRouter.Get("/statuses", orderHandler.GetOrderStatuses)
	orderRouter.Group("/:order_id").
		Post("/process", roleMiddleware.ValidateRole(user_service.UserRoleAdmin), orderHandler.OrderProcessing).
		Post("/out-for-delivery", roleMiddleware.ValidateRole(user_service.UserRoleAdmin), orderHandler.OutForDelivery).
		Post("/delivered", roleMiddleware.ValidateRole(user_service.UserRoleCustomer), orderHandler.Delivered)
}
