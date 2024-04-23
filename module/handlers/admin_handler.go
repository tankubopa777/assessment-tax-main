package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tankubopa777/assessment-tax/module/repository"
)

type AdminHandler struct {
	adminRepo repository.AdminRepository
}

func NewAdminHandler(adminRepo repository.AdminRepository) *AdminHandler {
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
		Deduction float64 `json:"personalDeduction"`
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	err := h.adminRepo.SetPersonalDeduction(request.Deduction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating personal deduction")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Personal deduction updated successfully",
	})
}

func (h *AdminHandler) SetKReceiptLimit(c echo.Context) error {
	var request struct {
		Limit float64 `json:"kReceiptLimit"`
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	err := h.adminRepo.SetKReceiptLimit(request.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating K-receipt limit")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "K-receipt limit updated successfully",
	})
}
