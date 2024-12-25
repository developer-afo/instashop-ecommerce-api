package router

import (
	"github.com/gofiber/fiber/v2"

	core_handler "github.com/developer-afo/instashop-ecommerce-api/handler/core"
	"github.com/developer-afo/instashop-ecommerce-api/lib/config"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/middleware"
	core_repository "github.com/developer-afo/instashop-ecommerce-api/repository/core"
	user_repository "github.com/developer-afo/instashop-ecommerce-api/repository/user"
	core_service "github.com/developer-afo/instashop-ecommerce-api/service/core"
	userService "github.com/developer-afo/instashop-ecommerce-api/service/user"
)

func InitializeCoreRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)
	productRepository := core_repository.NewProductRepository(db)
	imageRepository := core_repository.NewImageRepository(db)

	// Services
	imageService := core_service.NewImageService(imageRepository)
	productService := core_service.NewProductService(
		productRepository,
		imageService,
	)

	// config
	mediaConfig := config.NewMediaHelper(env)

	// Handlers
	productHandler := core_handler.NewProductHandler(productService, imageService)
	mediaHandler := core_handler.NewMediaHandler(mediaConfig)

	// middlewares
	authMiddleware := middleware.Protected()
	roleMiddleware := middleware.NewRoleMiddleware(userRepository)

	// Base routes
	productRoute := router.Group("/products")
	mediaRouter := router.Group("/media")

	// Routes

	productRoute.Post("/", authMiddleware, roleMiddleware.ValidateRole(userService.UserRoleAdmin), productHandler.CreateProduct)
	productRoute.Get("/", productHandler.FindAllProducts)
	productRoute.Get("/:slug", productHandler.FindProduct)
	productRoute.Put("/:id", authMiddleware, roleMiddleware.ValidateRole(userService.UserRoleAdmin), productHandler.UpdateProduct)
	productRoute.Delete("/:id", authMiddleware, roleMiddleware.ValidateRole(userService.UserRoleAdmin), productHandler.DeleteProduct)
	productRoute.Group("/:product_id/images", authMiddleware, roleMiddleware.ValidateRole(userService.UserRoleAdmin)).
		Get("/", productHandler.FindImagesByProductId).
		Post("/", productHandler.CreateImage).
		Delete("/:key", productHandler.DeleteImage)

	mediaRouter.Post("/upload", mediaHandler.UploadMedia, authMiddleware, roleMiddleware.ValidateRole(userService.UserRoleAdmin))
	mediaRouter.Get("/:key", mediaHandler.GetMedia)
}
