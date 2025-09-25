# AI Backend Service

A Go-based microservice for face recognition and user authentication, built with Gin framework and integrated with face recognition AI engine.

## Features

- User authentication (login/register)
- Face enrollment (register face identity)
- Face recognition (identify faces)
- Face deletion (remove face identity)
- JWT-based authentication
- gRPC integration
- RESTful API endpoints

## Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Authentication**: JWT
- **Database**: MySQL with GORM
- **Communication**: gRPC, REST API
- **Logging**: Structured logging with slog
- **Configuration**: Environment variables

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd ai-backend-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
# Copy and configure environment variables
cp .env.example .env
```

4. Run the application:
```bash
go run main.go
```

## Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `FACE_HOST` | Face recognition service host | `http://face-reg-engine:8080/face/v1/api/` |
| `FACE_ENROLL_API` | Face enrollment endpoint | `register-identity` |
| `FACE_RECOGNIZE_API` | Face recognition endpoint | `recognize-identity` |
| `FACE_DELETE_API` | Face deletion endpoint | `delete-identity` |
| `FACE_LIST_API` | Face list endpoint | `get-list` |
| `FACE_SERVICE_LOG_LEVEL` | Logging level | `info` |

## API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication

Most face-related endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <jwt-token>
```

---

## User Authentication APIs

### 1. User Registration
**Endpoint**: `POST /api/v1/user/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response**:
```json
{
  "data": true,
  "message": "Registration successful"
}
```

### 2. User Login
**Endpoint**: `POST /api/v1/user/authenticate`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_token_here",
    "expires_in": 3600
  }
}
```

---

## Face Recognition APIs

### 3. Face Enrollment (Register Face Identity)
**Endpoint**: `POST /api/v1/face/register-identity`

**Headers**:
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "userId": "user123",
  "imageBase64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ...",
  "gender": 1.0
}
```

**Response**:
```json
{
  "data": {
    "userId": "user123",
    "code": "200",
    "message": "Face enrolled successfully",
    "requestId": "req-12345",
    "data": {
      "name": "user123",
      "created_at": "2025-09-25T10:30:00Z",
      "image": "processed_image_base64"
    }
  }
}
```

### 4. Face Recognition (Identify Face)
**Endpoint**: `POST /api/v1/face/recognize-identity`

**Headers**:
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "imageBase64": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ..."
}
```

**Response**:
```json
{
  "data": {
    "userId": "user123",
    "code": "200",
    "message": "Face recognized successfully",
    "requestId": "req-12346",
    "data": {
      "name": "user123",
      "probability": 0.95,
      "created_at": "2025-09-25T10:31:00Z",
      "image": "processed_image_base64"
    }
  }
}
```

### 5. Face Deletion (Remove Face Identity)
**Endpoint**: `POST /api/v1/face/delete-identity`

**Headers**:
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "userId": "user123"
}
```

**Response**:
```json
{
  "data": {
    "userId": "user123",
    "code": "200",
    "message": "Face deleted successfully",
    "requestId": "req-12347",
    "data": {
      "name": "user123",
      "created_at": "2025-09-25T10:32:00Z"
    }
  }
}
```

---

## Profile APIs

### 6. Get User Profile
**Endpoint**: `POST /profile`

**Headers**:
```
Authorization: Bearer <jwt-token>
```

**Response**:
```json
{
  "data": {
    "id": "user123",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-09-25T09:00:00Z"
  }
}
```

---

## Health Check

### Ping
**Endpoint**: `GET /ping`

**Response**:
```json
{
  "data": "pong"
}
```

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": {
    "code": "400",
    "message": "Bad Request",
    "details": "Invalid request parameters"
  }
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid request parameters |
| 401 | Unauthorized - Missing or invalid JWT token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Development

### Project Structure

```
├── cmd/                    # Command line interface
├── common/                 # Common utilities and constants
├── composer/              # Service composition
├── config/                # Configuration management
├── controllers/           # HTTP controllers
├── helpers/               # Helper utilities
├── logger/                # Logging utilities
├── middleware/            # HTTP middleware
├── models/                # Data models
├── proto/                 # Protocol buffer definitions
├── service/               # Business logic services
├── transport/             # Transport layer (API handlers)
└── utils/                 # Utility functions
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o app main.go
```

### Docker

```bash
# Build image
docker build -t ai-backend-service .

# Run container
docker run -p 8080:8080 ai-backend-service
```

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
