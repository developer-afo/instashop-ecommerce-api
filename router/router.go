package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/developer-afo/instashop-ecommerce-api/handler"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
)

func InitializeRouter(router *fiber.App, dbConn database.DatabaseInterface, env constants.Env) {

	router.Get("/monitor", monitor.New(monitor.Config{Title: "Instashop API Monitor"}))

	InitializeUserRouter(router, dbConn, env)
	InitializeCoreRouter(router, dbConn, env)
	InitializeOrderRouter(router, dbConn, env)

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	router.Get("/", handler.Index)
	router.Get("*", handler.NotFound)

}
