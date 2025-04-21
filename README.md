# Accountability App

A video call application with accountability features, built with Go and Gin.

## Features

- RESTful API with Swagger/OpenAPI documentation
- Real-time WebSocket communication
- JWT-based authentication
- PostgreSQL database with GORM
- Structured error handling
- API versioning
- Comprehensive API documentation

## Prerequisites

- Go 1.24 or later
- PostgreSQL
- Swagger/OpenAPI tools

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/accountability-app.git
cd accountability-app
```

2. Install required tools:
```bash
make tools
```

3. Configure the application (edit config.yaml):
```yaml
database:
  host: localhost
  port: 5432
  user: your_user
  password: your_password
  name: videocall
  sslmode: disable

jwt:
  secret: your-secret-key
  expiration_hours: 24

websocket:
  allowed_origins:
    - http://localhost:3000
    - http://localhost:5173

server:
  port: 8080
  environment: development
```

4. Generate Swagger documentation:
```bash
make gen
```

5. Run the server:
```bash
make run
```

## API Documentation

### Swagger/OpenAPI Documentation
After starting the server, access the interactive API documentation at:
- http://localhost:8080/swagger/index.html

### Available Endpoints

#### User Service
- `POST /api/users/register` - Register a new user
  - Request: `CreateUserRequest` (email, username, password)
  - Response: `LoginResponse` (token, user details)

- `POST /api/users/login` - Authenticate user
  - Request: `LoginRequest` (email, password)
  - Response: `LoginResponse` (token, user details)

- `GET /api/users/{id}` - Get user details
  - Response: `UserResponse` (user details)
  - Requires: JWT Authentication

#### Call Service
- `POST /api/calls` - Create a new call
  - Request: `CallCreate` (title, description, creator_id)
  - Response: `Call` object

- `POST /api/calls/join` - Join an existing call
  - Request: `CallJoin` (call_id, user_id)
  - Response: `CallParticipant` object

- `POST /api/calls/{room_id}/leave` - Leave a call
  - Response: Success message
  - Requires: JWT Authentication

- `GET /api/calls` - List active calls
  - Response: Array of `Call` objects
  - Requires: JWT Authentication

#### WebSocket Service
- `GET /api/ws` - WebSocket connection endpoint
  - Query Parameters: room_id, user_id
  - Upgrades to WebSocket connection
  - Requires: JWT Authentication

- `GET /api/rooms/{room_id}/participants` - Get room participants count
  - Response: `ParticipantsResponse` (count)
  - Requires: JWT Authentication

## Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/          # Application entry point
│   ├── docs/                # Generated Swagger docs
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routes
│   │   ├── auth/           # Authentication logic
│   │   ├── config/         # Configuration management
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Database models
│   │   └── websocket/      # WebSocket implementation
│   ├── config.yaml         # Application configuration
│   └── go.mod             # Go module definition
└── Makefile              # Build and development commands
```

## Development Commands

- `make tools` - Install development tools
- `make gen` - Generate Swagger documentation
- `make clean` - Clean generated files
- `make run` - Run the server
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linter
- `make build` - Build the server binary

## Error Handling

The application uses a standardized error response format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## WebSocket Protocol

WebSocket connections require:
1. Valid JWT token in Authorization header
2. room_id query parameter
3. user_id query parameter

Messages are JSON-encoded with the following structure:
```json
{
  "type": "message_type",
  "payload": {}
}
```

## License

MIT 