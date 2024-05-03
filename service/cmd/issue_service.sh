read RA_WEBS_ADMIN_TOKEN
TTP_URL=http://localhost:8000
export RA_WEBS_SERVICE_TOKEN=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $RA_WEBS_ADMIN_TOKEN" $TTP_URL/service)