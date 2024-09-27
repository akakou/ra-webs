VERIFIER_BASE="https://tells-object-keep-threat.trycloudflare.com"
ADMIN_TOKEN="5af5745994ff116304f96a76bb05460bc7f80cd2dc0a6e5532ee085261b06576cac5dcabb86ce9026a80bdfd16bb0991009a564f4e8c044abde5f9c9f68c6868"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" $VERIFIER_BASE/api/notify -d "{\"domain\" : \"$1\", \"body\" : \"$2\"}"
