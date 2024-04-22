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
        {LowerBound: 0, UpperBound: 150000, Rate: 0},
        {LowerBound: 150000, UpperBound: 500000, Rate: 0.1},
        {LowerBound: 500000, UpperBound: 1000000, Rate: 0.15},
        {LowerBound: 1000000, UpperBound: 2000000, Rate: 0.2},
        {LowerBound: 2000000, UpperBound: -1, Rate: 0.35},
    }

    var totalDeductions float64 = 60000
    var donationDeduction float64 = 0
    var kReceiptDeduction float64 = 0

    for _, allowance := range input.Allowances {
        switch allowance.AllowanceType {
        case "donation":
            donationDeduction += allowance.Amount
            if donationDeduction > 100000 {
                donationDeduction = 100000
            }
        case "k-receipt":
            kReceiptDeduction += allowance.Amount
            if kReceiptDeduction > 50000 {
                kReceiptDeduction = 50000
            }
        }
    }

    totalDeductions += donationDeduction + kReceiptDeduction

    taxableIncome := input.TotalIncome - totalDeductions

    var totalTax float64
    taxDetails := make([]models.TaxLevelDetail, len(taxBrackets))

    for i, bracket := range taxBrackets {
        upperLimit := float64(bracket.UpperBound)
        if bracket.UpperBound == -1 {
            upperLimit = taxableIncome
        } else if taxableIncome < upperLimit {
            upperLimit = taxableIncome
        }

        lowerBound := float64(bracket.LowerBound)
        incomeInBracket := upperLimit - lowerBound
        taxForBracket := 0.0

        if taxableIncome > lowerBound && incomeInBracket > 0 {
            taxForBracket = incomeInBracket * bracket.Rate
            totalTax += taxForBracket
        }

        taxDetails[i] = models.TaxLevelDetail{
            Level: bracket.String(),
            Tax:   taxForBracket,
        }
    }

    finalTax := totalTax - input.WHT
    var taxRefund float64
    if finalTax < 0 {
        taxRefund = -finalTax
        finalTax = 0
    }

    return models.TaxCalculationResult{
        Tax:             finalTax,
        TaxRefund:       taxRefund,
        TaxLevelDetails: taxDetails,
    }, nil
}

