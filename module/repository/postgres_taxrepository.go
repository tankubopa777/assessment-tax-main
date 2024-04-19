package repository

import (
	"database/sql"
	"fmt"

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

    var totalDeductions float64 = 60000 // Standard personal deduction

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
    var taxAmount float64
    taxDetails := make([]models.TaxLevelDetail, 0)

    for _, bracket := range taxBrackets {
        if taxableIncome <= 0 {
            break
        }

        incomeInThisBracket := taxableIncome
        if bracket.UpperBound != -1 && taxableIncome > float64(bracket.UpperBound-bracket.LowerBound) {
            incomeInThisBracket = float64(bracket.UpperBound - bracket.LowerBound)
        }

        taxInThisBracket := incomeInThisBracket * bracket.Rate
        taxAmount += taxInThisBracket
        taxableIncome -= incomeInThisBracket

        levelDescription := fmt.Sprintf("%d-%d", bracket.LowerBound, bracket.UpperBound)
        if bracket.UpperBound == -1 {
            levelDescription = fmt.Sprintf("%d ขึ้นไป", bracket.LowerBound)
        }

        taxDetails = append(taxDetails, models.TaxLevelDetail{
            Level: levelDescription,
            Tax:   taxInThisBracket,
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
