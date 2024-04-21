package repository

import (
	"database/sql"
	"fmt"
	"math"

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

    var totalDeductions float64 = 60000 // Standard personal deduction

    // Process allowances and cap them accordingly
    for _, allowance := range input.Allowances {
        if allowance.AllowanceType == "donation" {
            totalDeductions += math.Min(allowance.Amount, 100000)
        } else if allowance.AllowanceType == "k-receipt" {
            totalDeductions += math.Min(allowance.Amount, 50000)
        }
    }

    taxableIncome := input.TotalIncome - totalDeductions
    // fmt.Println("taxableIncome: ", taxableIncome)
    var taxAmount float64
    taxDetails := make([]models.TaxLevelDetail, len(taxBrackets))

    // Populate the tax details with level descriptions
    for i, bracket := range taxBrackets {
        if bracket.UpperBound == -1 {
            taxDetails[i].Level = fmt.Sprintf("%d ขึ้นไป", bracket.LowerBound)
        } else {
            taxDetails[i].Level = fmt.Sprintf("%d-%d", bracket.LowerBound, bracket.UpperBound)
        }
    }

    // Calculate tax for the applicable income ranges
    for i, bracket := range taxBrackets {
        if taxableIncome > float64(bracket.LowerBound) {
            incomeInBracket := taxableIncome
            if bracket.UpperBound != -1 && taxableIncome > float64(bracket.UpperBound) {
                incomeInBracket = float64(bracket.UpperBound) - (float64(bracket.LowerBound) - 1)
            } else {
                incomeInBracket -= (float64(bracket.LowerBound) - 1)
            }

            taxAmount = incomeInBracket * bracket.Rate
            taxAmount = math.Round(taxAmount*100) / 100
            taxDetails[i].Tax = math.Round(taxAmount)
            if bracket.UpperBound != -1 && taxableIncome <= float64(bracket.UpperBound) {
                break
            }
        }
    }

    // Adjust tax by withholding tax and handle tax refund scenario
    finalTax := taxAmount - input.WHT
    var taxRefund float64
    if finalTax < 0 {
        taxRefund = -finalTax
        finalTax = 0
    }

    return models.TaxCalculationResult{
        Tax:            finalTax,
        TaxRefund:      taxRefund,
    }, nil
}
type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}
