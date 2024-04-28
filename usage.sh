curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 100000.0
    }
  ]
}'


curl -X POST http://localhost:8080/tax/calculations \
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

curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 200000.0
    },
    {
      "allowanceType": "donation",
      "amount": 100000.0
    }
  ]
}'
