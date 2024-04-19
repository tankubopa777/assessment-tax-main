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
