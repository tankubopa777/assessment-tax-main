package main

import (
	"log"
	"net/http"

	"github.com/KKGo-Software-engineering/assessment-tax/module/config"
	"github.com/KKGo-Software-engineering/assessment-tax/module/handlers"
	"github.com/KKGo-Software-engineering/assessment-tax/module/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Establish the database connection outside the handler.
	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err) // Using log.Fatalf will terminate the program if there's an error
	}
	defer db.Close()

	// Test database connection handler
	e.GET("/testdb", func(c echo.Context) error {
		return c.String(http.StatusOK, "Connected to database successfully!")
	})

	// Initialize the repository with the database connection
	taxRepo := repository.NewPostgresTaxRepository(db)

	// Initialize the handler with the repository
	taxHandler := handlers.NewTaxHandler(taxRepo)

	// Define the POST route
	e.POST("/tax/calculations", taxHandler.CalculateTax)

	// Start the server
	e.Logger.Fatal(e.Start(":5050"))
}
