# ADMIN_TOKEN=""
# VERIFIER_BASE=""

DOMAIN=$1
MESSAGE=$2
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" $VERIFIER_BASE/api/notify -d "{\"domain\" : \"$DOMAIN\", \"body\" : \"$MESSAGE\"}"
