package handlers

import (
	"net/http"

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

	result, err := h.repo.CalculateTax(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error calculating tax")
	}

	return c.JSON(http.StatusOK, result)
}
