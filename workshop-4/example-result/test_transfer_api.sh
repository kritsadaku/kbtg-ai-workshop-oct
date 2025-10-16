#!/bin/bash

# Test script for Transfer API
BASE_URL="http://localhost:3000"

echo "Testing Transfer API..."
echo "================================"

echo "1. Health check"
curl -X GET "$BASE_URL/" | jq .
echo ""

echo "2. Check existing users"
curl -X GET "$BASE_URL/users" | jq .
echo ""

echo "3. Create transfer from user 1 to user 3 (500 points)"
TRANSFER_RESPONSE=$(curl -s -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 3,
    "amount": 500,
    "note": "ทดสอบการโอนแต้ม"
  }')
echo $TRANSFER_RESPONSE | jq .

# Extract idemKey for later use
IDEM_KEY=$(echo $TRANSFER_RESPONSE | jq -r '.transfer.idemKey')
echo "Extracted idemKey: $IDEM_KEY"
echo ""

echo "4. Check updated points"
echo "Sender (User 1):"
curl -X GET "$BASE_URL/users/1" | jq .
echo ""
echo "Receiver (User 3):"
curl -X GET "$BASE_URL/users/3" | jq .
echo ""

echo "5. Get transfer by idempotency key"
curl -X GET "$BASE_URL/transfers/$IDEM_KEY" | jq .
echo ""

echo "6. Get transfer history for user 1"
curl -X GET "$BASE_URL/transfers?userId=1" | jq .
echo ""

echo "7. Get transfer history for user 3"
curl -X GET "$BASE_URL/transfers?userId=3" | jq .
echo ""

echo "8. Test error cases:"

echo "8.1. Insufficient points"
curl -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 3,
    "amount": 50000
  }' | jq .
echo ""

echo "8.2. Transfer to self"
curl -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 1,
    "amount": 100
  }' | jq .
echo ""

echo "8.3. Invalid user ID"
curl -X POST "$BASE_URL/transfers" \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 999,
    "toUserId": 3,
    "amount": 100
  }' | jq .
echo ""

echo "8.4. Non-existent transfer"
curl -X GET "$BASE_URL/transfers/non-existent-key" | jq .
echo ""

echo "8.5. Missing userId parameter"
curl -X GET "$BASE_URL/transfers" | jq .
echo ""

echo "9. Test pagination"
curl -X GET "$BASE_URL/transfers?userId=1&page=1&pageSize=10" | jq .
echo ""

echo "Transfer API testing completed!"