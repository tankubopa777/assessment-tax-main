package repository

import (
	"database/sql"
	"fmt"

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
    var kReceiptLimit float64 = 0

    err := r.db.QueryRow("SELECT personal_deduction FROM admin_settings WHERE id = 1").Scan(&totalDeductions)
    if err != nil {
        return models.TaxCalculationResult{}, fmt.Errorf("failed to fetch base deduction: %w", err)
    }

    err = r.db.QueryRow("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").Scan(&kReceiptLimit)
    if err != nil {
        return models.TaxCalculationResult{}, fmt.Errorf("failed to fetch k-receipt limit: %w", err)
    }

    for _, allowance := range input.Allowances {
        switch allowance.AllowanceType {
        case "donation":
            donationDeduction += allowance.Amount
            if donationDeduction > 100000 {
                donationDeduction = 100000
            }
        case "k-receipt":
            kReceiptDeduction += allowance.Amount
            if kReceiptDeduction > kReceiptLimit {
                kReceiptDeduction = kReceiptLimit
            }
        }
    }

    totalDeductions += donationDeduction + kReceiptDeduction

    taxableIncome := input.TotalIncome - totalDeductions

    var totalTax float64
    taxDetails := make([]models.TaxLevelDetail, len(taxBrackets))
    lastTaxedIndex := -1

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
            lastTaxedIndex = i
        }

        taxDetails[i] = models.TaxLevelDetail{
            Level: bracket.String(),
            Tax:   taxForBracket,
        }
    }

    if lastTaxedIndex != -1 {
        adjustedTax := taxDetails[lastTaxedIndex].Tax - input.WHT
        if adjustedTax < 0 {
            adjustedTax = 0
        }
        taxDetails[lastTaxedIndex].Tax = adjustedTax
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