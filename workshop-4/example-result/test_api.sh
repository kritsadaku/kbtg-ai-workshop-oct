#!/bin/bash

# Test script for User CRUD API
BASE_URL="http://localhost:3000"

echo "Testing User CRUD API..."
echo "================================"

echo "1. Health check"
curl -X GET "$BASE_URL/" | jq .
echo ""

echo "2. Get all users (should be empty initially)"
curl -X GET "$BASE_URL/users" | jq .
echo ""

echo "3. Create a new user (สมชาย ใจดี)"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "สมชาย",
    "last_name": "ใจดี", 
    "phone": "081-234-5678",
    "email": "somchai@example.com",
    "membership_level": "Gold",
    "points": 15420
  }' | jq .
echo ""

echo "4. Create another user"
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "สมหญิง",
    "last_name": "รักดี",
    "phone": "081-111-2222", 
    "email": "somying@example.com",
    "membership_level": "Silver",
    "points": 8500
  }' | jq .
echo ""

echo "5. Get all users"
curl -X GET "$BASE_URL/users" | jq .
echo ""

echo "6. Get user by ID (ID=1)"
curl -X GET "$BASE_URL/users/1" | jq .
echo ""

echo "7. Update user (ID=1)"
curl -X PUT "$BASE_URL/users/1" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 20000,
    "membership_level": "Platinum"
  }' | jq .
echo ""

echo "8. Get updated user (ID=1)"
curl -X GET "$BASE_URL/users/1" | jq .
echo ""

echo "9. Try to get non-existent user (ID=999)"
curl -X GET "$BASE_URL/users/999" | jq .
echo ""

echo "10. Delete user (ID=2)"
curl -X DELETE "$BASE_URL/users/2" | jq .
echo ""

echo "11. Get all users after deletion"
curl -X GET "$BASE_URL/users" | jq .
echo ""

echo "Testing completed!"