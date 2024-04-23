package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tankubopa777/assessment-tax/module/config"
	"github.com/tankubopa777/assessment-tax/module/handlers"
	"github.com/tankubopa777/assessment-tax/module/repository"
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	adminAuthMiddleware := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD"), nil
	})

	taxRepo := repository.NewPostgresTaxRepository(db)
    taxHandler := handlers.NewTaxHandler(taxRepo)
    adminRepo := repository.NewPostgresAdminRepository(db) 
    adminHandler := handlers.NewAdminHandler(adminRepo)    

    e.POST("/tax/calculations", taxHandler.CalculateTax)
    e.POST("/tax/calculations/upload-csv", taxHandler.UploadTaxCalculations) 

	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)

	adminGroup.GET("/settings", adminHandler.GetAdminSettings)
	adminGroup.POST("/deductions/personal", adminHandler.SetPersonalDeduction) 
	adminGroup.POST("/deductions/k-receipt-limit", adminHandler.SetKReceiptLimit) 

	e.Logger.Fatal(e.Start(":5050"))
}
