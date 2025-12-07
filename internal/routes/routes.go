package routes

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/kodra-pay/compliance-service/internal/handlers"
	"github.com/kodra-pay/compliance-service/internal/repositories"
	"github.com/kodra-pay/compliance-service/internal/services"
)

func Register(app *fiber.App, serviceName string) {
	// Health check
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	// Get database URL from environment
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		dbURL = "postgres://kodrapay:kodrapay_password@localhost:5432/kodrapay?sslmode=disable"
	} else {
		// Add sslmode=disable if not already present
		if !strings.Contains(dbURL, "sslmode=") {
			dbURL = dbURL + "?sslmode=disable"
		}
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize KYC components
	kycRepo := repositories.NewKYCRepository(db)
	kycService := services.NewKYCService(kycRepo)
	kycHandler := handlers.NewKYCHandler(kycService)

	// Register KYC routes
	kyc := app.Group("/kyc")
	kyc.Post("/submit", kycHandler.SubmitKYC)
	kyc.Get("/status/:merchant_id", kycHandler.GetKYCStatus)
	kyc.Post("/update", kycHandler.UpdateKYCStatus)
	kyc.Get("/pending", kycHandler.ListPendingKYC)
	kyc.Get("/list", kycHandler.ListKYCByStatus)
}
