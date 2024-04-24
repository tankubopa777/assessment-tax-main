package service

import (
	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
)

type TaxRepository interface {
	CalculateTax(input models.TaxCalculationInput) (models.TaxCalculationResult, error)
}

type AdminRepository interface {
}