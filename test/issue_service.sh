TTP_BASE="http://localhost:8000"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $1" $TTP_BASE/api/service
echo
