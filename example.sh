curl http://localhost:8080/1 -H 'X-API-KEY: SUPER_SECRET'

curl -X POST http://localhost:8080/1 \
    -d '{"key": "definition", "value": "420p"}' \
    -H 'X-API-KEY: SUPER_SECRET' \
    -H 'Content-Type: application/json' 

curl -X PUT POST http://localhost:8080/1 \
    -d '{"key": "definition", "value": "420p"}' \
    -H 'X-API-KEY: SUPER_SECRET' \
    -H 'Content-Type: application/json' 

curl -X DELETE POST http://localhost:8080/1 \
    -H 'X-API-KEY: SUPER_SECRET'
