package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/kodra-pay/compliance-service/internal/config"
	"github.com/kodra-pay/compliance-service/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Fiber app
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Register routes (includes DB wiring for KYC)
	routes.Register(app, "compliance-service")

	// Start server
	log.Printf("Compliance service starting on port %s", cfg.ServicePort)
	if err := app.Listen(":" + cfg.ServicePort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
