package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/compliance-service/internal/middleware"
	"github.com/kodra-pay/compliance-service/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7015"
	}

	serviceName := "compliance-service"

	app := fiber.New()
	app.Use(middleware.RequestID())

	routes.Register(app, serviceName)

	log.Printf("%s listening on :%s", serviceName, port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
