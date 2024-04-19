package models

type TaxBracket struct {
	LowerBound int    
	UpperBound int     
	Rate       float64 
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
	Tax         float64     
	TaxRefund   float64      
	TaxLevel    []TaxBracket 
}
