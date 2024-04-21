package repository

import (
	"database/sql"
	"testing"

	"github.com/KKGo-Software-engineering/assessment-tax/module/models"
	"github.com/stretchr/testify/assert"
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
			name: "Test Story 1",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       29000.0,
				TaxRefund: 0,
			},
			wantErr: false,
		},
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
				TaxRefund: 0,
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
				TaxRefund: 0,
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
				TaxRefund: 0,
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
				TaxRefund: 0,
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
				TaxRefund: 0,
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
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 1 : Tax Refund",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         40000.0,
				Allowances: []models.Allowance{
				
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       0.0,
				TaxRefund: 40000 - 29000, 
			},
			wantErr: false,
		},
		{
			name: "Test Case for donation allowance",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{
						AllowanceType: "donation",
						Amount:        200000.0,
					},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       19000,
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Case for k-receipt allowance",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []models.Allowance{
					{
						AllowanceType: "k-receipt",
						Amount:        60000.0,
					},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       24000,
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
		})
	}
}
