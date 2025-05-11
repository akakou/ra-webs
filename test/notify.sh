# RA_WEBS_SERVICE_TOKEN=""
# MONITOR_BASE=""

TITLE=$1
MESSAGE=$2
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $RA_WEBS_SERVICE_TOKEN" $MONITOR_BASE/api/notify -d "{\"message\" : \"$TITLE\n$MESSAGE\"}"
