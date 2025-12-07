package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kodra-pay/compliance-service/internal/models"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// ComplianceRepository defines the interface for compliance data operations
type ComplianceRepository interface {
	CreateKYCRecord(record *models.KYCRecord) error
	GetKYCRecordByID(id string) (*models.KYCRecord, error)
	UpdateKYCRecord(record *models.KYCRecord) error
	CreateTransactionMonitoringAlert(alert *models.TransactionMonitoringAlert) error
	GetTransactionMonitoringAlertByID(id string) (*models.TransactionMonitoringAlert, error)
	UpdateTransactionMonitoringAlert(alert *models.TransactionMonitoringAlert) error
}

// postgresComplianceRepository implements ComplianceRepository for PostgreSQL
type postgresComplianceRepository struct {
	db *sql.DB
}

// NewPostgresComplianceRepository creates a new PostgreSQL repository
func NewPostgresComplianceRepository(db *sql.DB) ComplianceRepository {
	return &postgresComplianceRepository{db: db}
}

func (r *postgresComplianceRepository) CreateKYCRecord(record *models.KYCRecord) error {
	query := `INSERT INTO kyc_records (id, user_id, status, document_type, document_id, issue_date, expiry_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, record.ID, record.UserID, record.Status, record.DocumentType, record.DocumentID, record.IssueDate, record.ExpiryDate, time.Now(), time.Now())
	return err
}

func (r *postgresComplianceRepository) GetKYCRecordByID(id string) (*models.KYCRecord, error) {
	record := &models.KYCRecord{}
	query := `SELECT id, user_id, status, document_type, document_id, issue_date, expiry_date, created_at, updated_at FROM kyc_records WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&record.ID, &record.UserID, &record.Status, &record.DocumentType, &record.DocumentID, &record.IssueDate, &record.ExpiryDate, &record.CreatedAt, &record.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Record not found
	}
	return record, err
}

func (r *postgresComplianceRepository) UpdateKYCRecord(record *models.KYCRecord) error {
	query := `UPDATE kyc_records SET user_id = $2, status = $3, document_type = $4, document_id = $5, issue_date = $6, expiry_date = $7, updated_at = $8 WHERE id = $1`
	_, err := r.db.Exec(query, record.ID, record.UserID, record.Status, record.DocumentType, record.DocumentID, record.IssueDate, record.ExpiryDate, time.Now())
	return err
}

func (r *postgresComplianceRepository) CreateTransactionMonitoringAlert(alert *models.TransactionMonitoringAlert) error {
	query := `INSERT INTO transaction_monitoring_alerts (id, transaction_id, user_id, rule_triggered, severity, status, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, alert.ID, alert.TransactionID, alert.UserID, alert.RuleTriggered, alert.Severity, alert.Status, alert.Description, time.Now(), time.Now())
	return err
}

func (r *postgresComplianceRepository) GetTransactionMonitoringAlertByID(id string) (*models.TransactionMonitoringAlert, error) {
	alert := &models.TransactionMonitoringAlert{}
	query := `SELECT id, transaction_id, user_id, rule_triggered, severity, status, description, created_at, updated_at FROM transaction_monitoring_alerts WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&alert.ID, &alert.TransactionID, &alert.UserID, &alert.RuleTriggered, &alert.Severity, &alert.Status, &alert.Description, &alert.CreatedAt, &alert.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Alert not found
	}
	return alert, err
}

func (r *postgresComplianceRepository) UpdateTransactionMonitoringAlert(alert *models.TransactionMonitoringAlert) error {
	query := `UPDATE transaction_monitoring_alerts SET transaction_id = $2, user_id = $3, rule_triggered = $4, severity = $5, status = $6, description = $7, updated_at = $8 WHERE id = $1`
	_, err := r.db.Exec(query, alert.ID, alert.TransactionID, alert.UserID, alert.RuleTriggered, alert.Severity, alert.Status, alert.Description, time.Now())
	return err
}

// InitDB initializes the database connection
func InitDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}
