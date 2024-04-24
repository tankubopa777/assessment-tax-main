package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
	"github.com/KKGo-Software-engineering/assessment-tax/module/repository"
	"github.com/labstack/echo/v4"
)

type TaxHandler struct {
	repo repository.TaxRepository
}

func NewTaxHandler(repo repository.TaxRepository) *TaxHandler {
	return &TaxHandler{
		repo: repo,
	}
}

func (h *TaxHandler) CalculateTax(c echo.Context) error {
	var input models.TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	result, err := h.repo.CalculateTax(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error calculating tax")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaxHandler) UploadTaxCalculations(c echo.Context) error {
	fileHeader, err := c.FormFile("taxFile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve the file from the form data.")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open the file.")
	}
	defer src.Close()

	
	tempDir := "uploads"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}

	tempFilePath := filepath.Join(tempDir, fileHeader.Filename)
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create a file for processing.")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy the file data.")
	}

	// Now that the file is saved, pass the file path to the repository method
	results, err := h.repo.TaxCalculationsFromCSV(tempFilePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to process the CSV file: %v", err))
	}

	// Optionally delete the file after processing if it's no longer needed
	os.Remove(tempFilePath)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"taxes": results,
	})
}