#!/bin/bash
API_URL="https://tchr-voice-server.sunburst.app/api"

curl  -d '{"name": "Biotech"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Health Care"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Information Technology"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Renewable Energy"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Entertainment"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Real Estate"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Blockchain"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
curl  -d '{"name": "Financial Services"}' -H 'Content-Type: application/json' -H "Authorization: Bearer $1" -X POST "$API_URL/industries"
