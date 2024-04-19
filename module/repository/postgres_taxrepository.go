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
	// Define tax brackets based on your exercise requirements.
	taxBrackets := []models.TaxBracket{
		{LowerBound: 0, UpperBound: 150000, Rate: 0},
		{LowerBound: 150001, UpperBound: 500000, Rate: 0.1},
		{LowerBound: 500001, UpperBound: 1000000, Rate: 0.15},
		{LowerBound: 1000001, UpperBound: 2000000, Rate: 0.2},
		{LowerBound: 2000001, UpperBound: -1, Rate: 0.35},
	}

	// Calculate the total deductions.
	var totalDeductions float64 = 60000 // The basic personal deduction.

	// Assuming that the AllowanceType will strictly be "donation" or "k-receipt" for this exercise.
	for _, allowance := range input.Allowances {
		totalDeductions += allowance.Amount
	}

	// Calculate taxable income.
	taxableIncome := input.TotalIncome - totalDeductions

	// Initialize the result with an empty slice for TaxLevelDetails.
	result := models.TaxCalculationResult{
		TaxLevelDetails: []models.TaxLevelDetail{},
	}

	// Calculate the tax based on the taxable income and populate TaxLevelDetails.
	for _, bracket := range taxBrackets {
		taxInThisBracket := 0.0
		if taxableIncome > 0 {
			incomeInBracket := min(taxableIncome, float64(bracket.UpperBound-bracket.LowerBound))
			taxInThisBracket = incomeInBracket * bracket.Rate
			result.Tax += taxInThisBracket
			taxableIncome -= incomeInBracket
		}

		// Append the tax detail for this bracket to the result.
		result.TaxLevelDetails = append(result.TaxLevelDetails, models.TaxLevelDetail{
			Level: bracket.String(), // Implement a String() method on TaxBracket if not already present.
			Tax:   taxInThisBracket,
		})

		// Break out of the loop if we have already covered all the income.
		if bracket.UpperBound == -1{
			break
		}
	}

	// Subtract Withholding Tax (WHT) if applicable.
	result.Tax -= input.WHT
	if result.Tax < 0 {
		result.TaxRefund = -result.Tax // If the result is negative, it means a refund is due.
		result.Tax = 0
	}

	return result, nil
}

// Helper function to get the minimum of two values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}
