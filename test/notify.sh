TTP_BASE="http://localhost:8000"
ADMIN_TOKEN="5e3e410b035405a78ff7b40724f914cf8df49ddd6bf6f1554f11274a066ae8f37917e131800cdf6869e7c659b95bc2d5007bf9b5456bf36c072b3d7a88d1e952"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" $TTP_BASE/api/notify -d '{"domain": "example.com", "body": "Hello, World!"}'

