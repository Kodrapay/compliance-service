package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/kodra-pay/compliance-service/internal/config"
	"github.com/kodra-pay/compliance-service/internal/handlers"
	"github.com/kodra-pay/compliance-service/internal/repositories"
	"github.com/kodra-pay/compliance-service/internal/routes"
	"github.com/kodra-pay/compliance-service/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	db, err := repositories.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup dependencies
	complianceRepo := repositories.NewPostgresComplianceRepository(db)
	complianceService := services.NewComplianceService(complianceRepo)
	complianceHandler := handlers.NewComplianceHandler(complianceService)

	// Setup routes
	routes.SetupComplianceRoutes(app, complianceHandler)

	// Start server
	log.Printf("Compliance service starting on port %s", cfg.ServicePort)
	if err := app.Listen(":" + cfg.ServicePort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
