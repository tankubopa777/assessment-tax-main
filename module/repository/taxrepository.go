package repository

import (
	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
)

type TaxRepository interface {
	CalculateTax(input models.TaxCalculationInput) (models.TaxCalculationResult, error)
	TaxCalculationsFromCSV(file string) ([]models.CSVTaxCalculationResult, error)
}

type AdminRepository interface {
	GetAdminSettings() (models.AdminSettings, error)
	SetPersonalDeduction(deduction float64) error
	SetKReceiptLimit(limit float64) error
}