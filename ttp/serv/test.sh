echo "Type admin token: "
read ADMIN_TOKEN

SERVICE_TOKEN=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN"  http://localhost:8081/service)
echo "Service Token is " $SERVICE_TOKEN

UNIQUE_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $SERVICE_TOKEN" -d '{"repository": "https://github.com/akakou-docs/ego-statistical-analysis"}' http://localhost:8081/code)
echo "Unique ID is " $UNIQUE_ID

SERVER_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $SERVICE_TOKEN" -d "{\"domain\": \"localhost:8080\" }" http://localhost:8081/server)
echo "Server ID is " $SERVER_ID
