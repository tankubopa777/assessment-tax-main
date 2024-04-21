package repository

import (
	"database/sql"

	"github.com/tankubopa777/assessment-tax/module/models"
)

type PostgresTaxRepository struct {
	db *sql.DB
}

func NewPostgresTaxRepository(db *sql.DB) *PostgresTaxRepository {
	return &PostgresTaxRepository{db: db}
}

func (r *PostgresTaxRepository) CalculateTax(input models.TaxCalculationInput) (models.TaxCalculationResult, error) {
    taxBrackets := []models.TaxBracket{
        {LowerBound: 150000, UpperBound: 500000, Rate: 0.1},
        {LowerBound: 500000, UpperBound: 1000000, Rate: 0.15},
        {LowerBound: 1000000, UpperBound: 2000000, Rate: 0.2},
        {LowerBound: 2000000, UpperBound: -1, Rate: 0.35},
    }

    var totalDeductions float64 = 60000 
    taxableIncome := input.TotalIncome - totalDeductions

    var totalTax float64

    for _, bracket := range taxBrackets {
        if taxableIncome > float64(bracket.LowerBound) {
            upperLimit := float64(bracket.UpperBound)
            if bracket.UpperBound == -1 {
                upperLimit = taxableIncome
            } else if taxableIncome < upperLimit {
                upperLimit = taxableIncome
            }

            incomeInBracket := upperLimit - float64(bracket.LowerBound)

            taxForBracket := incomeInBracket * bracket.Rate
            totalTax += taxForBracket
        }
    }

    finalTax := totalTax - input.WHT
    var taxRefund float64
    if finalTax < 0 {
        taxRefund = -finalTax
        finalTax = 0
    }

    return models.TaxCalculationResult{
        Tax:       finalTax,
        TaxRefund: taxRefund,
    }, nil
}
