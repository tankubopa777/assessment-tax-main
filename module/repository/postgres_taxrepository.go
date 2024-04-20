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
        {LowerBound: 2000001, UpperBound: -1, Rate: 0.35},
    }

    // Standard personal deduction
    var totalDeductions float64 = 60000 

    // Loop through allowances
    for _, allowance := range input.Allowances {
        switch allowance.AllowanceType {
        case "donation":
            if allowance.Amount > 100000 {
                totalDeductions += 100000
            } else {
                totalDeductions += allowance.Amount
            }
        case "k-receipt":
            if allowance.Amount > 50000 {
                totalDeductions += 50000
            } else {
                totalDeductions += allowance.Amount
            }
        }
    }

    taxableIncome := input.TotalIncome - totalDeductions
    fmt.Println("taxableIncome: ", taxableIncome)
    var taxAmount float64
    var applicableBracket *models.TaxBracket

    // Determine the tax bracket and calculate the tax based on the correct interval
    for _, bracket := range taxBrackets {
        if taxableIncome > float64(bracket.LowerBound) && (bracket.UpperBound == -1 || taxableIncome <= float64(bracket.UpperBound)) {
            applicableBracket = &bracket
            taxAmount = (taxableIncome - float64(bracket.LowerBound)) * bracket.Rate
            break
        }
    }

    taxAmount = math.Round(taxAmount)  // Round the tax amount to the nearest whole number

    var taxDetails []models.TaxLevelDetail
    if applicableBracket != nil {
        levelDescription := fmt.Sprintf("%d-%d", applicableBracket.LowerBound, applicableBracket.UpperBound)
        if applicableBracket.UpperBound == -1 {
            levelDescription = fmt.Sprintf("%d ขึ้นไป", applicableBracket.LowerBound)
        }
        taxDetails = append(taxDetails, models.TaxLevelDetail{
            Level: levelDescription,
            Tax:   taxAmount,
        })
    }

    finalTax := taxAmount - input.WHT
    var taxRefund float64
    if finalTax < 0 {
        taxRefund = -finalTax
        finalTax = 0
    }

    return models.TaxCalculationResult{
        Tax:            finalTax,
        TaxRefund:      taxRefund,
        TaxLevelDetails: taxDetails,
    }, nil
}




type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}
