package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/lib/seed"
	"github.com/developer-afo/instashop-ecommerce-api/router"
)

func main() {
	app := fiber.New(fiber.Config{AppName: "Instashop v0.0.1"})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               50,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	}))

	// Get environment variables
	env := constants.GetEnv()

	// Start database connection
	dbConn := database.StartDatabaseClient(env)

	// Initialize router
	router.InitializeRouter(app, dbConn, env)

	// Migrate database
	database.Migrate(dbConn)

	// Seed database
	seed.NewSeeder(dbConn).Seed()

	log.Fatal(app.Listen("0.0.0.0:" + env.PORT))
}
