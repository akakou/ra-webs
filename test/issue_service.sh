ADMIN_TOKEN=""
TTP_BASE="http://localhost:8000"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $token" $TTP_BASE/service
