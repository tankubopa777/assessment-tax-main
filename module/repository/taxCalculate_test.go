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
			name: "Test Kbank",
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
		{
			name: "Test for income 500000 with WHT 25000",
			input: models.TaxCalculationInput{
				TotalIncome: 500000.0,
				WHT:         25000.0,
				Allowances:  []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       4000,
				TaxRefund: 0,
				TaxLevelDetails: []models.TaxLevelDetail{
					{
						Level: "1-150000",
						Tax:   0,
					},
					{
						Level: "150001-500000",
						Tax:   4000,
					},
					{
						Level: "500001-1000000",
						Tax:   0,
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
		{
			name: "Test for income 500000 with WHT 30000",
			input: models.TaxCalculationInput{
				TotalIncome: 500000,
				WHT:         30000.0,
				Allowances:  []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       0,
				TaxRefund: 1000,
				TaxLevelDetails: []models.TaxLevelDetail{
					{
						Level: "1-150000",
						Tax:   0,
					},
					{
						Level: "150001-500000",
						Tax:   0,
					},
					{
						Level: "500001-1000000",
						Tax:   0,
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
		{
			name: "Test for income 500000 with donation 100000 and k-receipt 50000", 
			input: models.TaxCalculationInput{
				TotalIncome: 500000,
				WHT:         0.0,
				Allowances:  []models.Allowance{
					{
						AllowanceType: "donation",
						Amount:        100000,
					},
					{
						AllowanceType: "k-receipt",
						Amount:        50000,
					},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       14000,
				TaxRefund: 0,
				TaxLevelDetails: []models.TaxLevelDetail{
					{
						Level: "1-150000",
						Tax:   0,
					},
					{
						Level: "150001-500000",
						Tax:   14000,
					},
					{
						Level: "500001-1000000",
						Tax:   0,
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
			assert.Equal(t, tt.wantResult.Tax, gotResult.Tax, "Total calculated tax does not match expected tax")
			assert.Equal(t, tt.wantResult.TaxRefund, gotResult.TaxRefund, "Tax refund does not match expected refund")

			assert.Equal(t, len(tt.wantResult.TaxLevelDetails), len(gotResult.TaxLevelDetails), "Number of tax level details does not match")

			for i, wantDetail := range tt.wantResult.TaxLevelDetails {
    			if i < len(gotResult.TaxLevelDetails) { 
        			gotDetail := gotResult.TaxLevelDetails[i]
        			assert.Equal(t, wantDetail.Level, gotDetail.Level, "Tax level does not match at index", i)
        			assert.Equal(t, wantDetail.Tax, gotDetail.Tax, "Tax amount does not match at level", wantDetail.Level)
    			}
			}

		})
	}
}
