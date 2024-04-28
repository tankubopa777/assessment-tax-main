package service

import (
	"github.com/tankubopa777/assessment-tax/module/models"
)

type TaxRepository interface {
	CalculateTax(input models.TaxCalculationInput) (models.TaxCalculationResult, error)
}
