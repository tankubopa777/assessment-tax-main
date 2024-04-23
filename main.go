package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tankubopa777/assessment-tax/module/config"
	"github.com/tankubopa777/assessment-tax/module/handlers"
	"github.com/tankubopa777/assessment-tax/module/repository"
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

	go func() {
		if err := e.Start(":5050"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("shutting down the server")
}
