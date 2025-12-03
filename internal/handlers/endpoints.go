package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/compliance-service/internal/dto"
	"github.com/kodra-pay/compliance-service/internal/services"
)

type ComplianceHandler struct {
	svc *services.ComplianceService
}

func NewComplianceHandler(svc *services.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{svc: svc}
}

func (h *ComplianceHandler) WriteAudit(c *fiber.Ctx) error {
	var req dto.AuditLogRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.WriteAudit(c.Context(), req))
}

func (h *ComplianceHandler) ListAudit(c *fiber.Ctx) error {
	return c.JSON(h.svc.ListAudit(c.Context()))
}
