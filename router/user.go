package router

import (
	"github.com/gofiber/fiber/v2"

	userHandler "github.com/developer-afo/instashop-ecommerce-api/handler/user"
	"github.com/developer-afo/instashop-ecommerce-api/lib/config"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	user_repository "github.com/developer-afo/instashop-ecommerce-api/repository/user"
	"github.com/developer-afo/instashop-ecommerce-api/service"
	user_service "github.com/developer-afo/instashop-ecommerce-api/service/user"
)

func InitializeUserRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)
	verificationCodeRepository := user_repository.NewVerificationCodeRepository(db)

	// config
	mailConfig := config.NewEmail(env)

	// Services
	emailService := service.NewEmailService(mailConfig)
	userService := user_service.NewUserService(userRepository)
	verificationCodeService := user_service.NewVerficationCodeService(userRepository, verificationCodeRepository)
	authService := user_service.NewAuthService(userService, verificationCodeService, emailService)

	// Handler
	authHandler := userHandler.NewAuthHandler(authService)

	// Routers
	authRoute := router.Group("/auth")

	// Routes
	authRoute.Post("/login", authHandler.Login)
	authRoute.Post("/register", authHandler.Register)
	authRoute.Post("/refresh-token", authHandler.RefreshAccessToken)
	authRoute.Post("/resend-email", authHandler.ResendEmailVerification)
	authRoute.Post("/verify-email", authHandler.VerifyEmail)
	authRoute.Post("/forgot-password", authHandler.ForgotPassword)
	authRoute.Post("/reset-password", authHandler.ResetPassword)

}
