# Test for Kbank Tax Calculation API endpoint
curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,    
  "allowances": [{}]
}'

# Test for Kbank Tax Calculation API endpoint Request with WHT
curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 40000.0,    
  "allowances": [{}]
}'



