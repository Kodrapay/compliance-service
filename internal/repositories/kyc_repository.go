package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kodra-pay/compliance-service/internal/models"
)

type KYCRepository struct {
	db *sql.DB
}

func NewKYCRepository(db *sql.DB) *KYCRepository {
	return &KYCRepository{db: db}
}

// Create creates a new KYC submission
func (r *KYCRepository) Create(ctx context.Context, submission *models.KYCSubmission) error {
	// Convert documents map to JSONB
	docsJSON, err := json.Marshal(submission.Documents)
	if err != nil {
		return fmt.Errorf("failed to marshal documents: %w", err)
	}

	query := `
		INSERT INTO kyc_submissions (
			merchant_id, business_type, business_name, cac_number, tin_number,
			business_address, city, state, postal_code, incorporation_date,
			business_category, director_name, director_bvn, director_phone,
			director_email, documents, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 'pending')
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		submission.MerchantID,
		submission.BusinessType,
		submission.BusinessName,
		submission.CACNumber,
		submission.TINNumber,
		submission.BusinessAddress,
		submission.City,
		submission.State,
		submission.PostalCode,
		submission.IncorporationDate,
		submission.BusinessCategory,
		submission.DirectorName,
		submission.DirectorBVN,
		submission.DirectorPhone,
		submission.DirectorEmail,
		docsJSON,
	).Scan(&submission.ID, &submission.CreatedAt, &submission.UpdatedAt)
}

func (r *KYCRepository) GetByID(ctx context.Context, id int) (*models.KYCSubmission, error) {
	query := `
		SELECT id, merchant_id, business_type, business_name, cac_number, tin_number,
			business_address, city, state, postal_code, incorporation_date,
			business_category, director_name, director_bvn, director_phone,
			director_email, documents, status, reviewer_id, review_notes,
			reviewed_at, created_at, updated_at
		FROM kyc_submissions
		WHERE id = $1
	`

	var submission models.KYCSubmission
	var docsJSON []byte
	var reviewerID sql.NullInt32 // To handle nullable int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&submission.ID,
		&submission.MerchantID,
		&submission.BusinessType,
		&submission.BusinessName,
		&submission.CACNumber,
		&submission.TINNumber,
		&submission.BusinessAddress,
		&submission.City,
		&submission.State,
		&submission.PostalCode,
		&submission.IncorporationDate,
		&submission.BusinessCategory,
		&submission.DirectorName,
		&submission.DirectorBVN,
		&submission.DirectorPhone,
		&submission.DirectorEmail,
		&docsJSON,
		&submission.Status,
		&reviewerID,
		&submission.ReviewNotes,
		&submission.ReviewedAt,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal documents
	if err := json.Unmarshal(docsJSON, &submission.Documents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	// Handle nullable reviewerID
	if reviewerID.Valid {
		val := int(reviewerID.Int32)
		submission.ReviewerID = &val
	} else {
		submission.ReviewerID = nil
	}

	return &submission, nil
}

// GetLatestByMerchant retrieves the latest KYC submission for a merchant
func (r *KYCRepository) GetLatestByMerchant(ctx context.Context, merchantID int) (*models.KYCSubmission, error) {
	query := `
		SELECT id, merchant_id, business_type, business_name, cac_number, tin_number,
			business_address, city, state, postal_code, incorporation_date,
			business_category, director_name, director_bvn, director_phone,
			director_email, documents, status, reviewer_id, review_notes,
			reviewed_at, created_at, updated_at
		FROM kyc_submissions
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var submission models.KYCSubmission
	var docsJSON []byte
	var reviewerID sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, merchantID).Scan(
		&submission.ID,
		&submission.MerchantID,
		&submission.BusinessType,
		&submission.BusinessName,
		&submission.CACNumber,
		&submission.TINNumber,
		&submission.BusinessAddress,
		&submission.City,
		&submission.State,
		&submission.PostalCode,
		&submission.IncorporationDate,
		&submission.BusinessCategory,
		&submission.DirectorName,
		&submission.DirectorBVN,
		&submission.DirectorPhone,
		&submission.DirectorEmail,
		&docsJSON,
		&submission.Status,
		&reviewerID,
		&submission.ReviewNotes,
		&submission.ReviewedAt,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal documents
	if err := json.Unmarshal(docsJSON, &submission.Documents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	// Handle nullable reviewerID
	if reviewerID.Valid {
		val := int(reviewerID.Int32)
		submission.ReviewerID = &val
	} else {
		submission.ReviewerID = nil
	}

	return &submission, nil
}

// UpdateStatus updates the status of a KYC submission
func (r *KYCRepository) UpdateStatus(ctx context.Context, id int, status string, reviewerID *int, notes *string) error {
	query := `
		UPDATE kyc_submissions
		SET status = $1, reviewer_id = $2, review_notes = $3, reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, status, reviewerID, notes, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("kyc submission not found")
	}

	return nil
}

// ListByStatus retrieves KYC submissions by status
func (r *KYCRepository) ListByStatus(ctx context.Context, status string, limit int) ([]models.KYCSubmission, error) {
	query := `
		SELECT id, merchant_id, business_type, business_name, cac_number, tin_number,
			business_address, city, state, postal_code, incorporation_date,
			business_category, director_name, director_bvn, director_phone,
			director_email, documents, status, reviewer_id, review_notes,
			reviewed_at, created_at, updated_at
		FROM kyc_submissions
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []models.KYCSubmission
	for rows.Next() {
		var submission models.KYCSubmission
		var docsJSON []byte
		var reviewerID sql.NullInt32

		if err := rows.Scan(
			&submission.ID,
			&submission.MerchantID,
			&submission.BusinessType,
			&submission.BusinessName,
			&submission.CACNumber,
			&submission.TINNumber,
			&submission.BusinessAddress,
			&submission.City,
			&submission.State,
			&submission.PostalCode,
			&submission.IncorporationDate,
			&submission.BusinessCategory,
			&submission.DirectorName,
			&submission.DirectorBVN,
			&submission.DirectorPhone,
			&submission.DirectorEmail,
			&docsJSON,
			&submission.Status,
			&reviewerID,
			&submission.ReviewNotes,
			&submission.ReviewedAt,
			&submission.CreatedAt,
			&submission.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Unmarshal documents
		if err := json.Unmarshal(docsJSON, &submission.Documents); err != nil {
			return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
		}

		// Handle nullable reviewerID
		if reviewerID.Valid {
			val := int(reviewerID.Int32)
			submission.ReviewerID = &val
		} else {
			submission.ReviewerID = nil
		}

		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}
