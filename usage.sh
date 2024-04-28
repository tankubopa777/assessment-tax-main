# Admin get settings
curl -u adminTax:admin! http://localhost:8080/admin/settings

# Admin set settings
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":60000.0}' http://localhost:8080/admin/deductions/personal

# Admin set k
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":50000.0}' http://localhost:8080/admin/deductions/k-receipt-limit

# User test
curl -X POST http://localhost:8080/tax/calculations \
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
curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 500000.0,
  "wht": 0.0,
  "allowances": [{}]
}'

curl -X POST http://localhost:8080/tax/calculations \
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

curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 2160001.0,
  "wht": 200000.35,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 2000000.0
    }
  ]
}'

curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 7555533.0,
  "wht": 30000.0,
  "allowances": [
    {
      "allowanceType": "k-receipt",
      "amount": 0.0
    }
  ]
}'


curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 60000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 0.0
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
      "amount": 100000.0
    },
    {
      "allowanceType": "donation",
      "amount": 200000.0
    }
  ]
}'


600000,40000,20000
curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 600000.0,
  "wht": 40000.0,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 20000.0
    }
  ]
}'

curl -X POST http://localhost:8080/tax/calculations \
-H "Content-Type: application/json" \
-d '{
  "totalIncome": 600000.0,
  "wht": 0.0,
  "allowances": [
    {
      "allowanceType": "donation",
      "amount": 20000.0
    }
  ]
}'