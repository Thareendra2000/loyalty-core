# Loyalty Core API Test Commands

## Prerequisites
Make sure your server is running:
```bash
go run cmd/main.go
```

## 1. Health Check
```bash
curl http://localhost:8080/health
```

## 2. API Info
```bash
curl http://localhost:8080/api/info
```

## 3. List All Endpoints
```bash
curl http://localhost:8080/api/endpoints
```

## 4. Authentication Tests

### Sign Up
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 5. Loyalty Tests (All require authentication token)

### Earn Points
```bash
curl -X POST http://localhost:8080/api/loyalty/earn \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "points": 100,
    "description": "Purchase reward"
  }'
```

### Redeem Points
```bash
curl -X POST http://localhost:8080/api/loyalty/redeem \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "points": 50,
    "description": "Discount applied"
  }'
```

### Get Balance
```bash
curl -X GET http://localhost:8080/api/loyalty/balance \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get Transaction History
```bash
curl -X GET http://localhost:8080/api/loyalty/history \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get Transaction History with Limit
```bash
curl -X GET "http://localhost:8080/api/loyalty/history?limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 6. Test Flow Example

1. First, sign up a user
2. Login to get a token
3. Copy the token from the login response
4. Use the token in all subsequent requests
5. Test earning points, redeeming points, and checking balance

## Error Cases to Test

### Invalid Authentication
```bash
curl -X GET http://localhost:8080/api/loyalty/balance \
  -H "Authorization: Bearer invalid_token"
```

### Missing Authentication
```bash
curl -X GET http://localhost:8080/api/loyalty/balance
```

### Invalid Points (negative)
```bash
curl -X POST http://localhost:8080/api/loyalty/earn \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "points": -10,
    "description": "Invalid points"
  }'
```

### Insufficient Points for Redemption
```bash
curl -X POST http://localhost:8080/api/loyalty/redeem \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "points": 1000,
    "description": "More points than available"
  }'
```


## Quick Test Flow
1. Run health check to ensure server is running
2. Create a user with signup
3. Login with the same credentials to get a token
4. Use the token to get user profile
5. Try error cases (wrong password, invalid token, etc.)

## Using jq for Pretty JSON Output
If you have `jq` installed, you can pipe the output for better formatting:
```bash
curl -s http://localhost:8080/health | jq .
```
