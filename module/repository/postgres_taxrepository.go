package repository

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

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

func (r *PostgresTaxRepository) UploadTaxCalculations (input models.TaxCalculationInput, result models.TaxCalculationResult) error {
    return nil
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

// ProcessTaxCalculationsFromCSV processes tax calculations from a CSV file.
func (r *PostgresTaxRepository) ProcessTaxCalculationsFromCSV(filePath string) ([]models.CSVTaxCalculationResult, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    csvReader := csv.NewReader(file)
    var results []models.CSVTaxCalculationResult

    // Read the header first to skip it but also to check if the CSV is correctly formatted
    headers, err := csvReader.Read()
    if err != nil {
        return nil, err
    }
    if len(headers) < 3 || headers[0] != "totalIncome" || headers[1] != "wht" || headers[2] != "donation" {
        return nil, fmt.Errorf("CSV format is incorrect, expected headers: 'totalIncome, wht, donation'")
    }

    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        totalIncome, err := strconv.ParseFloat(record[0], 64)
        if err != nil {
            return nil, fmt.Errorf("error parsing totalIncome: %v", err)
        }
        wht, err := strconv.ParseFloat(record[1], 64)
        if err != nil {
            return nil, fmt.Errorf("error parsing wht: %v", err)
        }
        donation, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            return nil, fmt.Errorf("error parsing donation: %v", err)
        }

        // Calculate tax based on the parsed values
        taxableIncome := totalIncome - (60000 + donation) // Personal deduction and donation deduction
        tax := calculateTax(taxableIncome - wht) // Adjust taxable income by WHT and calculate tax

        results = append(results, models.CSVTaxCalculationResult{
            TotalIncome: totalIncome,
            Tax:         tax,
        })
    }

    return results, nil
}


func calculateTax(taxableIncome float64) float64 {
	taxBrackets := []models.TaxBracket{
		{LowerBound: 0, UpperBound: 150000, Rate: 0},
		{LowerBound: 150001, UpperBound: 500000, Rate: 0.1},
		{LowerBound: 500001, UpperBound: 1000000, Rate: 0.15},
		{LowerBound: 1000001, UpperBound: 2000000, Rate: 0.2},
		{LowerBound: 2000001, UpperBound: -1, Rate: 0.35},
	}

	var tax float64
	for _, bracket := range taxBrackets {
		if taxableIncome <= float64(bracket.LowerBound) {
			continue
		}
		// ตรวจสอบว่า upperBound ของช่วงภาษีนั้นเป็น -1 หรือไม่
		upperBound := float64(bracket.UpperBound)
		if bracket.UpperBound == -1 {
			upperBound = taxableIncome
		}

		// คำนวณรายได้ในช่วงนี้ที่ต้องเสียภาษี
		if taxableIncome > upperBound {
			tax += (upperBound - float64(bracket.LowerBound)) * bracket.Rate
		} else {
			tax += (taxableIncome - float64(bracket.LowerBound)) * bracket.Rate
			break
		}
	}
	return tax
}

