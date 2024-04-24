package repository

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tankubopa777/assessment-tax/module/models"
)

func TestCalculateTax(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresTaxRepository(db)

	tests := []struct {
		name        string
		mockQueries func()
		input       models.TaxCalculationInput
		wantResult  models.TaxCalculationResult
		wantErr     bool
	}{
		{
			name: "Basic Income Tax Calculation",
			mockQueries: func() {
				mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(60000))
				mock.ExpectQuery("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"k_receipt_limit"}).AddRow(100000))
			},
			input: models.TaxCalculationInput{
				TotalIncome: 500000,
				WHT:         0,
				Allowances: []models.Allowance{
					{AllowanceType: "donation", Amount: 200000},
					{AllowanceType: "k-receipt", Amount: 0},
				},
			},
			wantResult: models.TaxCalculationResult{
				Tax:       19000,
				TaxRefund: 0,
				TaxLevelDetails: []models.TaxLevelDetail{
					{Level: "0-150000", Tax: 0},
					{Level:"150001-500000", Tax: 19000},
					{Level:"500001-1000000", Tax: 0},
					{Level:"1000001-2000000",Tax: 0},
					{Level:"2000001 ขึ้นไป",Tax: 0},
				},
			},
			wantErr: false,
		},
		{
			name: "Handling Errors from Database",
			mockQueries: func() {
				mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
					WillReturnError(fmt.Errorf("database error"))
			},
			input: models.TaxCalculationInput{
				TotalIncome: 500000,
				WHT:         25000,
				Allowances:  []models.Allowance{},
			},
			wantResult: models.TaxCalculationResult{},
			wantErr:    true,
		},
		{
            name: "Error Fetching k_receipt_limit",
            mockQueries: func() {
                mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(60000))
                mock.ExpectQuery("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnError(fmt.Errorf("database error"))
            },
            input: models.TaxCalculationInput{TotalIncome: 100000},
            wantResult: models.TaxCalculationResult{},
            wantErr:    true,
        },
		{
            name: "k_receipt Deduction Greater Than Limit",
            mockQueries: func() {
                mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(60000))
                mock.ExpectQuery("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"k_receipt_limit"}).AddRow(10000))
            },
            input: models.TaxCalculationInput{
                TotalIncome: 500000,
                Allowances: []models.Allowance{
                    {AllowanceType: "k-receipt", Amount: 15000},
                },
            },
            wantResult: models.TaxCalculationResult{
                Tax: 28000,
                TaxRefund: 0,
                TaxLevelDetails: []models.TaxLevelDetail{
					{Level: "0-150000", Tax: 0},
					{Level:"150001-500000", Tax: 28000},
					{Level:"500001-1000000", Tax: 0},
					{Level:"1000001-2000000",Tax: 0},
					{Level:"2000001 ขึ้นไป",Tax: 0},
				},
            },
            wantErr: false,
        },
		{
            name: "WHT Greater Than Tax in Last Bracket",
            mockQueries: func() {
                mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(60000))
                mock.ExpectQuery("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"k_receipt_limit"}).AddRow(100000))
            },
            input: models.TaxCalculationInput{
                TotalIncome: 300000,  
                WHT:         50000,  
                Allowances:  []models.Allowance{},
            },
            wantResult: models.TaxCalculationResult{
                Tax: 0,  
                TaxRefund: 41000,
                TaxLevelDetails: []models.TaxLevelDetail{ 
					{Level: "0-150000", Tax: 0},
					{Level:"150001-500000", Tax: 0},
					{Level:"500001-1000000", Tax: 0},
					{Level:"1000001-2000000",Tax: 0},
					{Level:"2000001 ขึ้นไป",Tax: 0},
				 },
            },
            wantErr: false,
        },{
            name: "Total WHT Exceeds Total Calculated Tax",
            mockQueries: func() {
                mock.ExpectQuery("SELECT personal_deduction FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"personal_deduction"}).AddRow(60000))
                mock.ExpectQuery("SELECT k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnRows(sqlmock.NewRows([]string{"k_receipt_limit"}).AddRow(100000))
            },
            input: models.TaxCalculationInput{
                TotalIncome: 700000,  
                WHT:         120000,  
                Allowances:  []models.Allowance{},
            },
            wantResult: models.TaxCalculationResult{
                Tax: 0,
                TaxRefund: 64000,
                TaxLevelDetails: []models.TaxLevelDetail{ 
				{Level: "0-150000", Tax: 0},
				{Level:"150001-500000", Tax: 35000},
				{Level:"500001-1000000", Tax: 0},
				{Level:"1000001-2000000",Tax: 0},
				{Level:"2000001 ขึ้นไป",Tax: 0}, },
            },
            wantErr: false,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQueries() 
			gotResult, err := repo.CalculateTax(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresTaxRepository.CalculateTax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantResult.Tax, gotResult.Tax, "Calculated tax does not match expected tax")
				assert.Equal(t, tt.wantResult.TaxRefund, gotResult.TaxRefund, "Tax refund does not match expected refund")
				assert.Equal(t, len(tt.wantResult.TaxLevelDetails), len(gotResult.TaxLevelDetails), "Number of tax level details does not match")
				for i, wantDetail := range tt.wantResult.TaxLevelDetails {
					if i < len(gotResult.TaxLevelDetails) {
						gotDetail := gotResult.TaxLevelDetails[i]
						assert.Equal(t, wantDetail.Level, gotDetail.Level, "Tax level does not match at index", i)
						assert.Equal(t, wantDetail.Tax, gotDetail.Tax, "Tax amount does not match at level", wantDetail.Level)
					}
				}
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
