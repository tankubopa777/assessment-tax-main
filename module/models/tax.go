package models

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
	Tax         float64        `json:"tax,omitempty"`
	TaxRefund   float64        `json:"tax_refund,omitempty"`    
	TaxLevel    []TaxBracket   `json:"tax_level,omitempty"` 
}
