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

    var totalDeductions float64 = 60000 // ค่าลดหย่อนส่วนตัวเริ่มต้น
    var donationDeduction float64 = 0
    var kReceiptDeduction float64 = 0

    // Calculate total deductions
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
        } else {
            totalDeductions += allowance.Amount
        }
    }
    totalDeductions += donationDeduction + kReceiptDeduction

    taxableIncome := input.TotalIncome - totalDeductions

    result := models.TaxCalculationResult{
        TaxLevel: make([]models.TaxBracket, 0),
    }

    taxBeforeWHT := 0.0
    remainingIncome := taxableIncome

    for _, bracket := range taxBrackets {
        if remainingIncome > float64(bracket.LowerBound) {
            incomeInBracket := remainingIncome
            if bracket.UpperBound != -1 && remainingIncome > float64(bracket.UpperBound) {
                incomeInBracket = float64(bracket.UpperBound - bracket.LowerBound)
            }

            taxForBracket := incomeInBracket * bracket.Rate
            taxBeforeWHT += taxForBracket
            remainingIncome -= incomeInBracket

            result.TaxLevel = append(result.TaxLevel, models.TaxBracket{
                LowerBound: bracket.LowerBound,
                UpperBound: bracket.UpperBound,
                Rate:       bracket.Rate,
            })

            if bracket.UpperBound == -1 {
                break
            }
        }
    }

    taxAfterWHT := taxBeforeWHT - input.WHT
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