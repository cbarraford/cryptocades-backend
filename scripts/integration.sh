#!/bin/bash

set -x
set -e

curl -sX POST -d '{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345"}' "localhost:8080/users" | jq .
TOKEN=$(curl -sX POST -d '{"username":"bob","password":"password"}' "localhost:8080/login" | jq -r .token)
curl -s -H "Session: $TOKEN" "localhost:8080/users/1"
curl -sX DELETE -d "{\"token\":\"$TOKEN\"}" "localhost:8080/logout" | jq .
