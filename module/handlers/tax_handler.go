package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/tankubopa777/assessment-tax/module/models"
	"github.com/tankubopa777/assessment-tax/module/service"
)

type TaxHandler struct {
	repo service.TaxRepository
}

func NewTaxHandler(repo service.TaxRepository) *TaxHandler {
	return &TaxHandler{
		repo: repo,
	}
}

func (h *TaxHandler) CalculateTax(c echo.Context) error {
	var input models.TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}


	if input.WHT < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "WHT must be greater than or equal to 0")
	}

	if input.WHT > input.TotalIncome {
		return echo.NewHTTPError(http.StatusBadRequest, "WHT must be less than or equal to total income")
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

	results, err := h.repo.TaxCalculationsFromCSV(tempFilePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to process the CSV file: %v", err))
	}

	os.Remove(tempFilePath)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"taxes": results,
	})
}