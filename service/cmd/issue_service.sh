read RA_WEBS_ADMIN_TOKEN
RA_WEBS_SERVICE_TOKEN=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $RA_WEBS_ADMIN_TOKEN" http://localhost:8081/service)