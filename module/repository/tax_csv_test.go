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
		csvContent  string
		wantResults []models.CSVTaxCalculationResult
		wantErr     bool
	}{
		{
			name: "Valid CSV Data",
			csvContent: `totalIncome,wht,donation
500000,25000,20000
700000,50000,30000`,
			wantResults: []models.CSVTaxCalculationResult{
				{TotalIncome: 500000, Tax: 2000, TaxRefund: 0},
				{TotalIncome: 700000, Tax: 1500, TaxRefund: 0},
			},
			wantErr: false,
		},
		{
			name: "Invalid Headers",
			csvContent: `income,withholding,tax
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "*.csv")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name()) 

			if _, err := tempFile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			if err := tempFile.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			results, err := repo.TaxCalculationsFromCSV(tempFile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("Test '%s' failed: TaxCalculationsFromCSV() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantResults, results)
			}
		})
	}
}
