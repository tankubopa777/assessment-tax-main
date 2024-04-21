package repository

import (
	"database/sql"
	"fmt"
	"math"

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
        {LowerBound: 150001, UpperBound: 500000, Rate: 0.1},
        {LowerBound: 500001, UpperBound: 1000000, Rate: 0.15},
        {LowerBound: 1000001, UpperBound: 2000000, Rate: 0.2},
        {LowerBound: 2000001, UpperBound: -1, Rate: 0.35},
    }

    var totalDeductions float64 = 60000 

    for _, allowance := range input.Allowances {
        if allowance.AllowanceType == "donation" {
            totalDeductions += math.Min(allowance.Amount, 100000)
        } else if allowance.AllowanceType == "k-receipt" {
            totalDeductions += math.Min(allowance.Amount, 50000)
        }
    }

    taxableIncome := input.TotalIncome - totalDeductions
    fmt.Println("taxableIncome: ", taxableIncome)
    var taxAmount float64
    var taxDetails []models.TaxLevelDetail

    // Initialize all tax levels with zero tax
    for _, bracket := range taxBrackets {
        levelDescription := fmt.Sprintf("%d-%d", bracket.LowerBound, bracket.UpperBound)
        if bracket.UpperBound == -1 {
            levelDescription = fmt.Sprintf("%d ขึ้นไป", bracket.LowerBound)
        }
        taxDetails = append(taxDetails, models.TaxLevelDetail{
            Level: levelDescription,
            Tax:   0,
        })
    }

    // Calculate tax for the applicable bracket
    for i, bracket := range taxBrackets {
        if taxableIncome > float64(bracket.LowerBound) && (bracket.UpperBound == -1 || taxableIncome <= float64(bracket.UpperBound)) {
            taxAmount = (taxableIncome - (float64(bracket.LowerBound)-1)) * bracket.Rate
            fmt.Println(float64(bracket.LowerBound))
            fmt.Println("taxAmount: ", taxAmount)
            if input.WHT > 0 {
                taxDetails[i].Tax = math.Round(taxAmount) - input.WHT
                break
            }
            taxDetails[i].Tax = math.Round(taxAmount)
            break
        }
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

func (r *PostgresAdminRepository) GetAdminSettings() (models.AdminSettings, error) {
    settings := models.AdminSettings{}
    err := r.db.QueryRow("SELECT personal_deduction, k_receipt_limit FROM admin_settings WHERE id = 1").Scan(&settings.PersonalDeduction, &settings.KReceiptLimit)
    if err != nil {
        return settings, err
    }
    return settings, nil
}

func (r *PostgresAdminRepository) SetPersonalDeduction(deduction float64) error {
    _, err := r.db.Exec("UPDATE admin_settings SET personal_deduction = $1 WHERE id = 1", deduction)
    return err
}

func (r *PostgresAdminRepository) SetKReceiptLimit(limit float64) error {
    _, err := r.db.Exec("UPDATE admin_settings SET k_receipt_limit = $1 WHERE id = 1", limit)
    return err
}