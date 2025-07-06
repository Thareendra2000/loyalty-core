# Loyalty Core - Square-Integrated Backend

A Go backend for a loyalty program that integrates with the Square Loyalty API. This backend provides REST endpoints for user authentication, earning points, redeeming rewards, viewing balance, and transaction history.

## Features

- **User Authentication**: JWT-based signup and login
- **Square Integration**: Real-time integration with Square Loyalty API
- **Points Management**: Earn and redeem loyalty points
- **Transaction History**: View detailed transaction history from Square
- **Balance Tracking**: Real-time balance from Square Loyalty accounts
- **Automatic Account Creation**: Creates Square loyalty accounts for new users

## API Endpoints

### Authentication
- `POST /api/auth/signup` - User registration
- `POST /api/auth/login` - User login

### Loyalty Program (Requires Authentication)
- `POST /api/loyalty/earn` - Earn points
- `POST /api/loyalty/redeem` - Redeem points
- `GET /api/loyalty/balance` - Get current balance and recent transactions
- `GET /api/loyalty/history` - Get full transaction history

## Quick Start

### 1. Prerequisites
- Go 1.19 or higher
- Square Developer Account (for sandbox testing)

### 2. Environment Setup

Create a `.env` file in the root directory:

```env
PORT=8080
JWT_SECRET=your-very-secure-secret-key-here-change-this-in-production
SQUARE_ACCESS_TOKEN=your-square-sandbox-access-token
SQUARE_APPLICATION_ID=your-square-application-id
SQUARE_LOCATION_ID=your-square-location-id
SQUARE_ENVIRONMENT=sandbox
```

### 3. Square Setup

1. Create a Square Developer Account
2. Create a new application in the Square Dashboard
3. Subscribe to Square Loyalty (free in sandbox)
4. Get your credentials from the Square Dashboard:
   - Access Token (from OAuth or use sandbox token)
   - Application ID
   - Location ID (from Locations API or Dashboard)

### 4. Run the Application

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/main.go
```

The server will start on port 8080 (or your configured port).

## Testing the API

Use the provided test commands in `API_TEST_COMMANDS.md` or use the following examples:

### 1. User Registration
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

### 2. User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Earn Points (use JWT token from login)
```bash
curl -X POST http://localhost:8080/api/loyalty/earn \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "points": 100,
    "description": "Purchase reward"
  }'
```

### 4. Get Balance
```bash
curl -X GET http://localhost:8080/api/loyalty/balance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Project Structure

```
loyalty-core/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── controllers/
│   ├── auth_controller.go     # Authentication endpoints
│   └── loyalty_controller.go  # Loyalty program endpoints
├── middleware/
│   └── auth.go               # JWT authentication middleware
├── models/
│   ├── user.go               # User data models
│   └── transaction.go        # Transaction data models
├── routes/
│   ├── auth_routes.go        # Authentication routes
│   ├── loyalty_routes.go     # Loyalty program routes
│   └── main_router.go        # Main router setup
├── services/
│   ├── auth_service.go       # Authentication business logic
│   ├── loyalty_service.go    # Loyalty program business logic
│   └── square_service.go     # Square API integration
├── storage/
│   └── user_storage.go       # In-memory user storage
├── utils/
│   └── jwt.go                # JWT utilities
├── .env                      # Environment variables
├── .env.example             # Environment template
├── API_TEST_COMMANDS.md     # API testing commands
├── go.mod                   # Go module dependencies
└── README.md               # This file
```

## Square Integration Details

The application integrates with Square Loyalty API using the official Square Go SDK:

### Key Integration Points:
1. **Automatic Account Creation**: When a user earns points for the first time, a Square loyalty account is automatically created
2. **Real-time Points**: Points are tracked in real-time through Square's API
3. **Transaction History**: Transaction history is fetched from Square's loyalty events
4. **Balance Synchronization**: Balance is always fetched from Square for accuracy

### Square API Operations Used:
- `CreateLoyaltyAccount` - Create loyalty accounts for new users
- `AccumulateLoyaltyPoints` - Add points when users earn them
- `AdjustLoyaltyPoints` - Subtract points when users redeem them
- `GetLoyaltyAccount` - Get current balance and account info
- `SearchLoyaltyEvents` - Get transaction history

## Error Handling

The application includes comprehensive error handling:
- Invalid credentials return 401 Unauthorized
- Missing authentication returns 401 Unauthorized
- Insufficient points returns 400 Bad Request
- Square API errors are properly handled and logged
- Fallback to local storage if Square API is unavailable

## Security Features

- JWT-based authentication
- Password hashing using bcrypt
- Request validation and sanitization
- Environment-based configuration
- Secure token handling

## Next Steps

1. **Frontend Integration**: Ready for React frontend integration
2. **Database**: Replace in-memory storage with persistent database
3. **Webhooks**: Add Square webhook support for real-time updates
4. **Rewards Catalog**: Integrate with Square's reward tiers
5. **Analytics**: Add loyalty program analytics and reporting

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.