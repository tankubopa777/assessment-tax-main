package main

import (
	"log"
	"os"

	"github.com/KKGo-Software-engineering/assessment-tax/module/config"
	"github.com/KKGo-Software-engineering/assessment-tax/module/handlers"
	"github.com/KKGo-Software-engineering/assessment-tax/module/repository"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	
	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Middleware for logging and recovering from panics
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the Basic Auth middleware for admin routes
	adminAuthMiddleware := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD"), nil
	})

	taxRepo := repository.NewPostgresTaxRepository(db)
    taxHandler := handlers.NewTaxHandler(taxRepo)
    adminRepo := repository.NewPostgresAdminRepository(db) // Create admin repository instance
    adminHandler := handlers.NewAdminHandler(adminRepo)    // Create admin handler instance

    // Public endpoints
    e.POST("/tax/calculations", taxHandler.CalculateTax)
    e.POST("/tax/calculations/upload-csv", taxHandler.UploadTaxCalculations) // New endpoint for uploading CSV

	// Admin endpoints with Basic Auth Middleware
	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)

	// Admin setting routes using the adminHandler
	adminGroup.GET("/settings", adminHandler.GetAdminSettings) // หากต้องการเรียกดูการตั้งค่า
	adminGroup.POST("/deductions/personal", adminHandler.SetPersonalDeduction) // สำหรับการตั้งค่า personal deduction
	adminGroup.POST("/deductions/k-receipt-limit", adminHandler.SetKReceiptLimit) // สำหรับการตั้งค่า k-receipt limit

	e.Logger.Fatal(e.Start(":5050"))
}
