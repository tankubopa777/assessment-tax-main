package handlers

import (
	"net/http"

	"github.com/KKGo-Software-engineering/assessment-tax/module/repository"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminRepo repository.AdminRepository
}

func NewAdminHandler(adminRepo repository.AdminRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo: adminRepo,
	}
}

// อย่าลืม import ของแพคเกจที่จำเป็น เช่น "net/http" และ "github.com/labstack/echo/v4"

// GetAdminSettings returns the admin settings to the client.
func (h *AdminHandler) GetAdminSettings(c echo.Context) error {
	settings, err := h.adminRepo.GetAdminSettings()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving admin settings")
	}
	return c.JSON(http.StatusOK, settings)
}

// SetPersonalDeduction updates the personal deduction in the admin settings.
func (h *AdminHandler) SetPersonalDeduction(c echo.Context) error {
	var request struct {
		Deduction float64 `json:"amount"`
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

// SetKReceiptLimit updates the K-receipt limit in the admin settings.
func (h *AdminHandler) SetKReceiptLimit(c echo.Context) error {
	var request struct {
		Limit float64 `json:"amount"`
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
