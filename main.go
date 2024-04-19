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

	db, err := config.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	e.GET("/testdb", func(c echo.Context) error {
		return c.String(http.StatusOK, "Connected to database successfully!")
	})

	taxRepo := repository.NewPostgresTaxRepository(db)

	taxHandler := handlers.NewTaxHandler(taxRepo)

	e.POST("/tax/calculations", taxHandler.CalculateTax)

	e.Logger.Fatal(e.Start(":5050"))
}
