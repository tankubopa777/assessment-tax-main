package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tankubopa777/assessment-tax/module/models"
)

func TestTaxCalculationsFromCSV(t *testing.T) {
	repo := NewPostgresTaxRepository(nil)

	tests := []struct {
		name        string
		setup       func(t *testing.T) (fileName string, cleanup func())
		csvContent  string
		wantResults []models.CSVTaxCalculationResult
		wantErr     bool
	}{
		{
			name: "Valid CSV Data",
			csvContent: `totalIncome,wht,donation
500000,25000,20000
700000,50000,30000
5000000,500000,200000`,
			wantResults: []models.CSVTaxCalculationResult{
				{TotalIncome: 500000, Tax: 2000, TaxRefund: 0},
				{TotalIncome: 700000, Tax: 1500, TaxRefund: 0},
				{TotalIncome: 5000000, Tax: 769000, TaxRefund: 0},
			},
			wantErr: false,
		},
		{
			name: "Invalid Headers",
			csvContent: `
500000,25000,20000`,
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Invalid Data Format",
			csvContent: `totalIncome,wht,donation
wrong,50000,20000
600000,25000,error`,
			wantResults: nil,
			wantErr:     true,
		},
		{
            name: "File Open Error",
            csvContent: "",
            wantResults: nil,
            wantErr:     true,
        },
        {
            name: "Empty CSV Content",
            csvContent: "totalIncome,wht,donation\n",
            wantResults: []models.CSVTaxCalculationResult(nil),
            wantErr:     false,
        },
        {
            name: "Tax Refund Scenario",
            csvContent: `totalIncome,wht,donation
500000,55000,20000`,
            wantResults: []models.CSVTaxCalculationResult{
                {TotalIncome: 500000, Tax: 0, TaxRefund: 28000}, 
            },
            wantErr:     false,
        },
		{
			name: "Malformed CSV Record",
			csvContent: `totalIncome,wht,donation
500000,"not_a_number",20000`,
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error Parsing WHT",
			csvContent: `totalIncome,wht,donation
500000,not_a_number,20000`,
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error Parsing Donation",
			csvContent: `totalIncome,wht,donation
500000,25000,not_a_number`, 
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error Parsing Donation",
			csvContent: ``,
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "CSV Record Read Error",
			csvContent: `totalIncome,wht,donation
		"500000", "malformed_quoting_20000`, 
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error Reading Headers",
			setup: func(t *testing.T) (string, func()) {
				tempFile, err := os.CreateTemp("", "restricted_*.csv")
				assert.NoError(t, err)
				tempFile.Close()
				os.Chmod(tempFile.Name(), 0222) 
				return tempFile.Name(), func() {
					os.Remove(tempFile.Name())
				}
			},
			wantResults: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fileName string
			var cleanup func()
			if tt.setup != nil {
				fileName, cleanup = tt.setup(t)
				defer cleanup()
			} else {
				tempFile, err := os.CreateTemp("", "*.csv")
				assert.NoError(t, err)
				defer os.Remove(tempFile.Name())

				if _, err := tempFile.WriteString(tt.csvContent); err != nil {
					t.Fatalf("failed to write to temp file: %v", err)
				}
				if err := tempFile.Close(); err != nil {
					t.Fatalf("failed to close temp file: %v", err)
				}
				fileName = tempFile.Name()
			}

			results, err := repo.TaxCalculationsFromCSV(fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test '%s' failed: TaxCalculationsFromCSV() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.Equal(t, tt.wantResults, results)
			}
		})
	}
}