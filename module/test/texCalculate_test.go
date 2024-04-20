package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
	"github.com/KKGo-Software-engineering/assessment-tax/module/repository"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	db := &sql.DB{}
	repo := repository.NewPostgresTaxRepository(db)

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
				Tax:       29000.0,
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 500,001-1,000,000",
			input: models.TaxCalculationInput{
				TotalIncome: 1000000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       66000.0,
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 1,000,001-2,000,000",
			input: models.TaxCalculationInput{
				TotalIncome: 1500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       88000.0,
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax level 2,000,001 and above",
			input: models.TaxCalculationInput{
				TotalIncome: 2500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       154000.0,
				TaxRefund: 0,
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