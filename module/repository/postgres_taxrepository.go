package repository

import (
	"database/sql"

	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
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
        {LowerBound: 150001, UpperBound: 500000, Rate: 0.1},
        {LowerBound: 500001, UpperBound: 1000000, Rate: 0.15},
        {LowerBound: 1000001, UpperBound: 2000000, Rate: 0.2},
        {LowerBound: 2000001, UpperBound: -1, Rate: 0.35}, // -1 indicates no upper limit
    }

    // Calculate total deductions
    var totalDeductions float64 = 60000 // Standard personal deduction
    var donationDeduction float64 = 0
    var kReceiptDeduction float64 = 0

    for _, allowance := range input.Allowances {
        if allowance.AllowanceType == "donation" {
            donationDeduction += allowance.Amount
            if donationDeduction > 100000 {
                donationDeduction = 100000
            }
        } else if allowance.AllowanceType == "k-receipt" {
            kReceiptDeduction += allowance.Amount
            if kReceiptDeduction > 50000 {
                kReceiptDeduction = 50000
            }
        }
    }
    totalDeductions += donationDeduction + kReceiptDeduction

    taxableIncome := input.TotalIncome - totalDeductions
    taxBeforeWHT := 0.0
    remainingIncome := taxableIncome

    // Calculating tax without storing tax levels
    for _, bracket := range taxBrackets {
        if remainingIncome > float64(bracket.LowerBound) {
            incomeInBracket := remainingIncome
            if bracket.UpperBound != -1 && remainingIncome > float64(bracket.UpperBound) {
                incomeInBracket = float64(bracket.UpperBound) - float64(bracket.LowerBound)
            }

            taxBeforeWHT += incomeInBracket * bracket.Rate
            remainingIncome -= incomeInBracket

            if bracket.UpperBound == -1 {
                break
            }
        }
    }

    taxAfterWHT := taxBeforeWHT - input.WHT
    result := models.TaxCalculationResult{}
    if taxAfterWHT < 0 {
        result.TaxRefund = -taxAfterWHT
        result.Tax = 0
    } else {
        result.Tax = taxAfterWHT
    }

    return result, nil
}

type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}
