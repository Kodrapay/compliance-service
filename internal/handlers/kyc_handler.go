package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/compliance-service/internal/dto"
	"github.com/kodra-pay/compliance-service/internal/services"
)

type KYCHandler struct {
	service *services.KYCService
}

func NewKYCHandler(service *services.KYCService) *KYCHandler {
	return &KYCHandler{service: service}
}

// SubmitKYC handles KYC submission requests
func (h *KYCHandler) SubmitKYC(c *fiber.Ctx) error {
	var req dto.KYCSubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.service.Submit(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetKYCStatus retrieves the KYC status for a merchant
func (h *KYCHandler) GetKYCStatus(c *fiber.Ctx) error {
	merchantID := c.Params("merchant_id")
	if merchantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "merchant_id is required")
	}

	status, err := h.service.GetLatest(c.Context(), merchantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get KYC status")
	}

	if status == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"merchant_id": merchantID,
			"status":      "not_started",
			"message":     "No KYC submission found",
		})
	}

	return c.JSON(status)
}

// UpdateKYCStatus updates the KYC status (admin only)
func (h *KYCHandler) UpdateKYCStatus(c *fiber.Ctx) error {
	var req dto.KYCStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.MerchantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "merchant_id is required")
	}
	if req.Status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "status is required")
	}

	if err := h.service.UpdateStatus(c.Context(), req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"merchant_id":  req.MerchantID,
		"status":       req.Status,
		"reviewer_id":  req.ReviewerID,
		"review_notes": req.ReviewNotes,
		"message":      "KYC status updated successfully",
	})
}

// ListPendingKYC lists all pending KYC submissions
func (h *KYCHandler) ListPendingKYC(c *fiber.Ctx) error {
	limit := 100
	if limitParam := c.QueryInt("limit", 100); limitParam > 0 && limitParam <= 100 {
		limit = limitParam
	}

	result, err := h.service.ListByStatus(c.Context(), "pending", limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list pending KYC submissions")
	}

	return c.JSON(result)
}

// ListKYCByStatus lists KYC submissions by status
func (h *KYCHandler) ListKYCByStatus(c *fiber.Ctx) error {
	status := c.Query("status", "pending")
	limit := c.QueryInt("limit", 100)

	if limit <= 0 || limit > 100 {
		limit = 100
	}

	result, err := h.service.ListByStatus(c.Context(), status, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list KYC submissions")
	}

	return c.JSON(result)
}
