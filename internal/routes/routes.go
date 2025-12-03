package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/compliance-service/internal/handlers"
	"github.com/kodra-pay/compliance-service/internal/services"
)

func Register(app *fiber.App, service string) {
	health := handlers.NewHealthHandler(service)
	health.Register(app)

	svc := services.NewComplianceService()
	h := handlers.NewComplianceHandler(svc)
	api := app.Group("/audit")
	api.Post("/logs", h.WriteAudit)
	api.Get("/logs", h.ListAudit)
}
