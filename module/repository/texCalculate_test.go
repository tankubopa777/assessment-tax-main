package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tankubopa777/assessment-tax/module/models"
)

func TestCalculateTax(t *testing.T) {
	db := &sql.DB{}
	repo := NewPostgresTaxRepository(db)

	tests := []struct {
		name       string
		input      models.TaxCalculationInput
		wantResult models.TaxCalculationResult
		wantErr    bool
	}{
		{
			name: "Test Story 1 : KBank want",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       29000,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 0 - 150,000",
			input: models.TaxCalculationInput{
				TotalIncome: 15000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       0.0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 500,001-1,000,000",
			input: models.TaxCalculationInput{
				TotalIncome: 660000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:         50000, 
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 500,001-1,000,000",
			input: models.TaxCalculationInput{
				TotalIncome: 1000000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:         101000, 
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 1,000,001-2,000,000",
			input: models.TaxCalculationInput{
				TotalIncome: 1500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       198000,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 2,000,001 and above",
			input: models.TaxCalculationInput{
				TotalIncome: 2500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       464000,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax Refund",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         40000.0,
				Allowances: []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				TaxRefund: 40000 - 29000, 
			},
			wantErr: false,
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := repo.CalculateTax(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresTaxRepository.CalculateTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantResult.Tax, gotResult.Tax, "Calculated tax does not match expected tax")
			if tt.wantErr {
				fmt.Printf("%s: PASSED (expected error)\n", tt.name)
			} else {
				fmt.Printf("%s: PASSED (expected result)\n", tt.name)
			}
		})
	}
}