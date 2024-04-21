# Admin get settings
curl -u adminTax:admin! http://localhost:5050/admin/settings

# Admin set settings
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":60000.0}' http://localhost:5050/admin/deductions/personal

# Admin set k
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":50000.0}' http://localhost:5050/admin/deductions/k-receipt-limit

# User test
curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 50000.0 
    },
    {
      "allowanceType": "donation",
      "amount": 100000.0
    }
  ]
}'

# User test
curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [{}]
}'

curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 50000.0 
    },
    {
      "allowanceType": "donation",
      "amount": 200000.0
    }
  ]
}'

curl -X POST http://localhost:5050/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 100000.0
    }
  ]
}'