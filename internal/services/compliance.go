package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/kodra-pay/compliance-service/internal/dto"
)

type ComplianceService struct{}

func NewComplianceService() *ComplianceService { return &ComplianceService{} }

func (s *ComplianceService) WriteAudit(_ context.Context, req dto.AuditLogRequest) dto.AuditLogResponse {
	return dto.AuditLogResponse{ID: "audit_" + uuid.NewString()}
}

func (s *ComplianceService) ListAudit(_ context.Context) []dto.AuditLogResponse {
	return []dto.AuditLogResponse{}
}
