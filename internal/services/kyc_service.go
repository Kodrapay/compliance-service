package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kodra-pay/compliance-service/internal/dto"
	"github.com/kodra-pay/compliance-service/internal/models"
	"github.com/kodra-pay/compliance-service/internal/repositories"
)

type KYCService struct {
	repo *repositories.KYCRepository
}

func NewKYCService(repo *repositories.KYCRepository) *KYCService {
	return &KYCService{repo: repo}
}

// Submit processes a KYC submission request
func (s *KYCService) Submit(ctx context.Context, req dto.KYCSubmissionRequest) (*dto.KYCSubmissionResponse, error) {
	// Validate required fields
	if req.MerchantID == 0 { // int check
		return nil, fmt.Errorf("merchant_id is required")
	}
	if req.BusinessName == "" {
		return nil, fmt.Errorf("business_name is required")
	}

	// Normalize business type

businessType := strings.ToLower(strings.TrimSpace(req.BusinessType))
	if businessType == "" {
	
businessType = "registered"
	}
	if businessType != "registered" && businessType != "startup" {
		return nil, fmt.Errorf("business_type must be 'registered' or 'startup'")
	}

	// Create submission model
	submission := &models.KYCSubmission{
		MerchantID:       req.MerchantID, // int
		BusinessType:     businessType,
		BusinessName:     req.BusinessName,
		CACNumber:        req.CACNumber,
		TINNumber:        req.TINNumber,
		BusinessAddress:  req.BusinessAddress,
		City:             req.City,
		State:            req.State,
		PostalCode:       req.PostalCode,
		BusinessCategory: req.BusinessCategory,
		DirectorName:     req.DirectorName,
		DirectorBVN:      req.DirectorBVN,
		DirectorPhone:    req.DirectorPhone,
		DirectorEmail:    req.DirectorEmail,
		Documents:        req.Documents,
		Status:           "pending",
	}

	// Parse incorporation date if provided
	if req.IncorporationDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.IncorporationDate); err == nil {
			submission.IncorporationDate = &parsed
		}
	}

	// Save to database
	if err := s.repo.Create(ctx, submission); err != nil {
		return nil, fmt.Errorf("failed to create KYC submission: %w", err)
	}

	// Update merchant KYC status to pending via merchant service
	if err := s.updateMerchantKYCStatus(req.MerchantID, "pending"); err != nil { // int
		// Log error but don't fail the submission
		fmt.Printf("Warning: failed to update merchant KYC status: %v\n", err)
	}

	return &dto.KYCSubmissionResponse{
		SubmissionID: submission.ID, // int
		Status:       "pending",
		Message:      "KYC submission received and is under review",
	}, nil
}

// GetLatest retrieves the latest KYC submission for a merchant
func (s *KYCService) GetLatest(ctx context.Context, merchantID int) (*dto.KYCStatusResponse, error) { // int
	submission, err := s.repo.GetLatestByMerchant(ctx, merchantID) // int
	if err != nil {
		return nil, fmt.Errorf("failed to get KYC status: %w", err)
	}
	if submission == nil {
		return nil, nil
	}

	return &dto.KYCStatusResponse{
		MerchantID:  submission.MerchantID, // int
		Status:      submission.Status,
		SubmittedAt: submission.CreatedAt.Format(time.RFC3339),
		ReviewedAt:  timePtrToString(submission.ReviewedAt),
		ReviewerID:  intPtrToInt(submission.ReviewerID), // *int
		ReviewNotes: stringPtrToString(submission.ReviewNotes),
	}, nil
}

// UpdateStatus updates the KYC status (admin operation)
func (s *KYCService) UpdateStatus(ctx context.Context, req dto.KYCStatusUpdateRequest) error {
	// Validate status
	status := strings.ToLower(req.Status)
	if status != "approved" && status != "rejected" && status != "pending" {
		return fmt.Errorf("invalid status: must be 'approved', 'rejected', or 'pending'")
	}

	// Get the latest submission for this merchant
	latest, err := s.repo.GetLatestByMerchant(ctx, req.MerchantID) // int
	if err != nil || latest == nil {
		return fmt.Errorf("no KYC submission found for merchant")
	}

	// Update status in database
	var reviewerID *int // Now *int
	if req.ReviewerID != 0 {
		reviewerID = &req.ReviewerID
	}
	notes := &req.ReviewNotes
	if err := s.repo.UpdateStatus(ctx, latest.ID, status, reviewerID, notes); err != nil { // int, *int
		return fmt.Errorf("failed to update KYC status: %w", err)
	}

	// Sync merchant KYC status
	if err := s.updateMerchantKYCStatus(req.MerchantID, status); err != nil { // int
		// Log error but don't fail the update
		fmt.Printf("Warning: failed to sync merchant KYC status: %v\n", err)
	}

	return nil
}

// ListByStatus lists KYC submissions by status
func (s *KYCService) ListByStatus(ctx context.Context, status string, limit int) (*dto.KYCListResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	submissions, err := s.repo.ListByStatus(ctx, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list KYC submissions: %w", err)
	}

	responses := make([]dto.KYCStatusResponse, 0, len(submissions))
	for _, sub := range submissions {
		responses = append(responses, dto.KYCStatusResponse{
			MerchantID:  sub.MerchantID, // int
			Status:      sub.Status,
			SubmittedAt: sub.CreatedAt.Format(time.RFC3339),
			ReviewedAt:  timePtrToString(sub.ReviewedAt),
			ReviewerID:  intPtrToInt(sub.ReviewerID), // *int
			ReviewNotes: stringPtrToString(sub.ReviewNotes),
		})
	}

	return &dto.KYCListResponse{
		Submissions: responses,
		Total:       len(responses),
	}, nil
}

// Helper function to update merchant KYC status via merchant service
func (s *KYCService) updateMerchantKYCStatus(merchantID int, status string) error { // int
	// Call merchant service to update KYC status
	url := fmt.Sprintf("http://merchant-service:7002/merchants/%d/kyc-status", merchantID) // int, changed %s to %d

	payload := map[string]string{
		"kyc_status": status,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("merchant service returned status %d", resp.StatusCode)
	}

	return nil
}

// Helper functions
func timePtrToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// intPtrToInt converts *int to int, returning 0 if nil
func intPtrToInt(i *int) int {
	if i == nil {
		return 0 // Or handle as an error, depending on requirements
	}
	return *i
}