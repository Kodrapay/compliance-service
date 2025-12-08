package dto

// KYCSubmissionRequest represents the request to submit KYC
type KYCSubmissionRequest struct {
	MerchantID        int               `json:"merchant_id"`
	BusinessType      string            `json:"business_type"` // "registered" or "startup"
	BusinessName      string            `json:"business_name"`
	CACNumber         string            `json:"cac_number,omitempty"`
	TINNumber         string            `json:"tin_number,omitempty"`
	BusinessAddress   string            `json:"business_address"`
	City              string            `json:"city"`
	State             string            `json:"state"`
	PostalCode        string            `json:"postal_code,omitempty"`
	IncorporationDate string            `json:"incorporation_date,omitempty"` // ISO 8601 format
	BusinessCategory  string            `json:"business_category"`
	DirectorName      string            `json:"director_name"`
	DirectorBVN       string            `json:"director_bvn"`
	DirectorPhone     string            `json:"director_phone"`
	DirectorEmail     string            `json:"director_email"`
	Documents         map[string]string `json:"documents"` // document_type -> file_path/url
}

// KYCSubmissionResponse represents the response after KYC submission
type KYCSubmissionResponse struct {
	SubmissionID int `json:"submission_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// KYCStatusUpdateRequest represents a request to update KYC status (admin only)
type KYCStatusUpdateRequest struct {
	MerchantID  int    `json:"merchant_id"`
	Status      string `json:"status"` // "approved", "rejected", or "pending"
	ReviewerID  int    `json:"reviewer_id"`
	ReviewNotes string `json:"review_notes,omitempty"`
}

// KYCStatusResponse represents the KYC status for a merchant
type KYCStatusResponse struct {
	MerchantID  int    `json:"merchant_id"`
	Status      string `json:"status"`
	SubmittedAt string `json:"submitted_at,omitempty"`
	ReviewedAt  string `json:"reviewed_at,omitempty"`
	ReviewerID  int    `json:"reviewer_id,omitempty"`
	ReviewNotes string `json:"review_notes,omitempty"`
}

// KYCListResponse represents a list of KYC submissions
type KYCListResponse struct {
	Submissions []KYCStatusResponse `json:"submissions"`
	Total       int                 `json:"total"`
}
