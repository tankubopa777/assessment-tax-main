package models

import "fmt"

type TaxBracket struct {
	LowerBound int 
	UpperBound int     
	Rate       float64 `json:"tax,omitempty"` 
}

type Allowance struct {
	AllowanceType string
	Amount        float64
}

type TaxCalculationInput struct {
	TotalIncome float64
	WHT         float64
	Allowances  []Allowance
}

type TaxCalculationResult struct {
    Tax             float64 `json:"tax,omitempty"`         
    TaxRefund       float64 `json:"taxRefund,omitempty"`   
    TaxLevelDetails []TaxLevelDetail `json:"taxLevel,omitempty"`
}
type TaxLevelDetail struct {
	Level string `json:"level"`
	Tax   float64 `json:"tax"`
}

func (tb TaxBracket) String() string {
	if tb.UpperBound == 150000 {
		return fmt.Sprintf("%d-%d", tb.LowerBound, tb.UpperBound)
	}
    if tb.UpperBound == -1 {
        return fmt.Sprintf("%d ขึ้นไป", tb.LowerBound + 1)
    }
    return fmt.Sprintf("%d-%d", tb.LowerBound + 1, tb.UpperBound)
}

