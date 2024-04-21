# Admin get settings
curl -u adminTax:admin! http://localhost:5050/admin/settings

# Admin set settings
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":70000.0}' http://localhost:5050/admin/deductions/personal

# Admin set k
curl -u adminTax:admin! -X POST -H "Content-Type: application/json" -d '{"amount":70000.0}' http://localhost:5050/admin/deductions/k-receipt-limit


