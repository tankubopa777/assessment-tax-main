package repository

import (
	"database/sql"
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
			name: "Test for income 660000",
			input: models.TaxCalculationInput{
				TotalIncome: 660000.0,
				WHT:         0.0,
				Allowances:  []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       50000,
				TaxRefund: 0,
				TaxLevelDetails: []models.TaxLevelDetail{
					{
						Level: "1-150000",
						Tax:   0,
					},
					{
						Level: "150001-500000",
						Tax:   35000,
					},
					{
						Level: "500001-1000000",
						Tax:   15000,
					},
					{
						Level: "1000001-2000000",
						Tax:   0,
					},
					{
						Level: "2000001 ขึ้นไป",
						Tax:   0,
					},
				},
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
