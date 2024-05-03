echo 'enter admin token: '
read RA_WEBS_ADMIN_TOKEN

echo 'enter domain: '
read RA_WEBS_DOMAIN

TTP_URL=http://localhost:8000
export RA_WEBS_SERVICE_TOKEN=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $RA_WEBS_ADMIN_TOKEN" $TTP_URL/service)

go run main.go code https://github.com/akakou-docs/ego-statistical-analysis
go run main.go server $RA_WEBS_DOMAIN 
