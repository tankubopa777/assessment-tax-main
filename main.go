package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)

	adminGroup.GET("/settings", adminHandler.GetAdminSettings) 
	adminGroup.POST("/deductions/personal", adminHandler.SetPersonalDeduction) 
	adminGroup.POST("/deductions/k-receipt-limit", adminHandler.SetKReceiptLimit) 

	go func() {
		if err := e.Start(":5050"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("shutting down the server failed", err)
	}

	log.Println("Server gracefully shutdown")
}
