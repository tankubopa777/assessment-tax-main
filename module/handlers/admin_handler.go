package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tankubopa777/assessment-tax/module/service"
)

type AdminHandler struct {
	adminRepo service.AdminRepository
}

func NewAdminHandler(adminRepo service.AdminRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo: adminRepo,
	}
}

func (h *AdminHandler) GetAdminSettings(c echo.Context) error {
	settings, err := h.adminRepo.GetAdminSettings()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving admin settings")
	}
	return c.JSON(http.StatusOK, settings)
}

func (h *AdminHandler) SetPersonalDeduction(c echo.Context) error {
	var request struct {
		Deduction float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	if request.Deduction > 100000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Personal deduction cannot exceed 100,000 THB")
	}

	if request.Deduction < 10000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Personal deduction cannot be less than 10,000 THB")
	}

	err := h.adminRepo.SetPersonalDeduction(request.Deduction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating personal deduction")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"personalDeduction": request.Deduction,
	})
}

func (h *AdminHandler) SetKReceiptLimit(c echo.Context) error {
	var request struct {
		Limit float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	if request.Limit > 100000 {
		return echo.NewHTTPError(http.StatusBadRequest, "K-receipt limit cannot exceed 100,000 THB")
	}

	if request.Limit <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "K-receipt limit must be greater than 0")
	}

	err := h.adminRepo.SetKReceiptLimit(request.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating K-receipt limit")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"kReceipt": request.Limit,
	})
}
