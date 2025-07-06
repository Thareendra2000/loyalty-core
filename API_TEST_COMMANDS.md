# Loyalty Core API Test Commands

## Prerequisites
Make sure your server is running:
```bash
go run cmd/main.go
```

## Basic Endpoints

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. API Info
```bash
curl http://localhost:8080/api/info
```

### 3. Available Endpoints
```bash
curl http://localhost:8080/api/endpoints
```

## Authentication Endpoints

### 4. User Signup
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

### 5. User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 6. Get User Profile
**Note:** Replace `YOUR_TOKEN_HERE` with the actual token from login response
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Test Cases

### Test Invalid Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }'
```

### Test Duplicate Signup
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "firstName": "Jane",
    "lastName": "Smith"
  }'
```

### Test Profile Without Token
```bash
curl -X GET http://localhost:8080/api/auth/profile
```

### Test Profile With Invalid Token
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer invalid_token"
```

## Expected Responses

### Successful Signup Response:
```json
{
  "message": "User created successfully",
  "user": {
    "id": "user_id_here",
    "email": "test@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "loyaltyId": "LOY12345678",
    "points": 0,
    "createdAt": "2025-07-06T...",
    "updatedAt": "2025-07-06T..."
  }
}
```

### Successful Login Response:
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_id_here",
    "email": "test@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "loyaltyId": "LOY12345678",
    "points": 0,
    "createdAt": "2025-07-06T...",
    "updatedAt": "2025-07-06T..."
  }
}
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
