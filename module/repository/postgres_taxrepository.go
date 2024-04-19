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
    // Define the tax brackets based on your exercise
    taxBrackets := []models.TaxBracket{
        {LowerBound: 0, UpperBound: 150000, Rate: 0},
        {LowerBound: 150001, UpperBound: 500000, Rate: 0.1},
        {LowerBound: 500001, UpperBound: 1000000, Rate: 0.15},
        {LowerBound: 1000001, UpperBound: 2000000, Rate: 0.2},
        {LowerBound: 2000001, UpperBound: -1, Rate: 0.35}, // -1 indicates no upper limit
    }

    // Calculate the total deductions
    var totalDeductions float64 = 60000 // The basic personal deduction
    for _, allowance := range input.Allowances {
        totalDeductions += allowance.Amount
    }

    // Calculate taxable income
    taxableIncome := input.TotalIncome - totalDeductions

    // Initialize the result
    result := models.TaxCalculationResult{
        TaxLevel: make([]models.TaxBracket, 0),
    }

    // Calculate the tax based on the taxable income
    for _, bracket := range taxBrackets {
        if taxableIncome > float64(bracket.LowerBound) {
            incomeInBracket := taxableIncome
            if bracket.UpperBound != -1 && taxableIncome > float64(bracket.UpperBound) {
                incomeInBracket = float64(bracket.UpperBound - bracket.LowerBound)
            }

            taxForBracket := incomeInBracket * bracket.Rate
            result.Tax += taxForBracket

            result.TaxLevel = append(result.TaxLevel, models.TaxBracket{
                LowerBound: bracket.LowerBound,
                UpperBound: bracket.UpperBound,
                Rate:       bracket.Rate,
            })

            if bracket.UpperBound == -1 || taxableIncome <= float64(bracket.UpperBound) {
                break
            }
        }
    }

    result.Tax -= input.WHT
    if result.Tax < 0 {
        result.TaxRefund = -result.Tax // Negative tax means refund
        result.Tax = 0
    }

    return result, nil
}

type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}

func (r *PostgresAdminRepository) GetAdminSettings() (models.AdminSettings, error) {
	// Fetch the admin settings from the database
	settings := models.AdminSettings{}
	// Database operation...
	return settings, nil
}

func (r *PostgresAdminRepository) SetPersonalDeduction(deduction float64) error {
	// Update the personal deduction in the database
	// Database operation...
	return nil
}

func (r *PostgresAdminRepository) SetKReceiptLimit(limit float64) error {
	// Update the KReceipt limit in the database
	// Database operation...
	return nil
}