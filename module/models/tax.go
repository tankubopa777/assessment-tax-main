package models

import "fmt"

// TaxBracket defines the tax rate for a range of income.
type TaxBracket struct {
	LowerBound int
	UpperBound int
	Rate       float64
}

// Allowance defines a type of deduction and its amount.
type Allowance struct {
	AllowanceType string
	Amount        float64
}

// TaxCalculationInput represents the input data for tax calculation.
type TaxCalculationInput struct {
	TotalIncome float64
	WHT         float64
	Allowances  []Allowance
}

// TaxCalculationResult represents the result of tax calculation.
type TaxCalculationResult struct {
	Tax             float64 `json:"tax"`
	TaxRefund       float64
	TaxLevelDetails []TaxLevelDetail `json:"taxLevel"`
}

// TaxLevelDetail represents the details of tax calculation for a specific income bracket.
type TaxLevelDetail struct {
	Level string `json:"level"`
	Tax   float64 `json:"tax"`
}

// String provides a string representation of a TaxBracket, which can be used for output.
func (tb TaxBracket) String() string {
	if tb.UpperBound == -1 {
		return fmt.Sprintf("%d ขึ้นไป", tb.LowerBound)
	}
	return fmt.Sprintf("%d-%d", tb.LowerBound, tb.UpperBound)
}
