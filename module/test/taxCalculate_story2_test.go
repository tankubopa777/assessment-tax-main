package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tankubopa777/assessment-tax/module/models"
	"github.com/tankubopa777/assessment-tax/module/repository"
)

func TestCalculateTaxStory2(t *testing.T) {
	db := &sql.DB{}
	repo := repository.NewPostgresTaxRepository(db)

	tests := []struct {
		name       string
		input      models.TaxCalculationInput
		wantResult models.TaxCalculationResult
		wantErr    bool
	}{
		{
			name: "Test Story 2 : KBank want",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         25000.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:      4000,
				TaxRefund: 0,
			},
			wantErr: false,
		},
		{
			name: "Test Story 2 : Tax level 500,001-1,000,000 with WHT",
			input: models.TaxCalculationInput{
				TotalIncome: 1000000.0,
				WHT:         25000.0,
				Allowances: []models.Allowance{
					{},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:      (((1000000 - 60000)-500000) * 0.15)-25000,
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