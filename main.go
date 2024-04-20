package main

import (
	"log"
	"net/http"
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

	// Public endpoints
	e.POST("/tax/calculations", taxHandler.CalculateTax)

	// Admin endpoints with Basic Auth Middleware
	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)

	// Here you would add your admin routes, e.g., to set personal deduction:
	adminGroup.POST("/deductions/personal", func(c echo.Context) error {
		// Your logic to handle the admin deduction setting
		return c.JSON(http.StatusOK, map[string]interface{}{
			"personalDeduction": 70000.0,
		})
	})

	e.Logger.Fatal(e.Start(":5050"))
}
