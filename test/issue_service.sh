# ADMIN_TOKEN=""
# RA_WEBS_VERIFIER_BASES=""

export RA_WEBS_SERVICE_TOKEN=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" $RA_WEBS_VERIFIER_BASES/api/service)

echo -en $RA_WEBS_SERVICE_TOKEN

