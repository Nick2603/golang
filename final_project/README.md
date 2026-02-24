# Final Project — REST API

A RESTful backend API for a blogging platform built with Go, Fiber, MongoDB, and JWT authentication.

## Features

- User registration and login with JWT-based authentication
- Secure password hashing with bcrypt
- Full CRUD for posts with ownership authorization
- Interactive API documentation via Swagger UI
- Dockerized deployment with Docker Compose
- Structured JSON logging and graceful shutdown

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Web Framework | [Fiber v3](https://github.com/gofiber/fiber) |
| Database | MongoDB |
| Authentication | JWT (`golang-jwt/jwt/v5`) |
| Password Hashing | bcrypt (`golang.org/x/crypto`) |
| API Docs | Swagger (`swaggo/swag`) |
| Containerization | Docker + Docker Compose |

## Project Structure

```
final_project/
├── cmd/server/
│   ├── main.go              # Entry point, server setup, routing
│   ├── handlers/            # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── post_handler.go
│   ├── middlewares/
│   │   └── auth.go          # JWT authentication middleware
│   └── utils/
│       └── response.go      # Response helpers
├── internal/
│   ├── clients/
│   │   └── mongodb.go       # MongoDB connection
│   ├── config/
│   │   └── config.go        # Configuration from environment
│   ├── models/
│   │   ├── user.go
│   │   └── post.go
│   ├── services/
│   │   ├── user_service.go  # User business logic
│   │   └── post_service.go  # Post business logic
│   └── utils/
│       ├── jwt.go           # JWT helpers
│       └── password.go      # Password helpers
├── docs/                    # Auto-generated Swagger docs
├── .env.example             # Environment variable template
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

## Prerequisites

- [Go 1.25+](https://golang.org/dl/)
- [MongoDB](https://www.mongodb.com/try/download/community) (local) **or** Docker

## Getting Started

### 1. Configure environment

```bash
cp .env.example .env
```

Edit `.env` as needed:

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=backend
JWT_SECRET=change-this-to-a-secure-random-secret
JWT_EXPIRY_HOURS=24
```

### 2. Run locally

```bash
# Install swag CLI (first time only)
make install-swag

# Generate Swagger docs and start the server
make run
```

### 3. Run with Docker Compose

```bash
# Start app + MongoDB containers
make docker-up

# Tail application logs
make docker-logs

# Stop and remove containers
make docker-down
```

### 4. Build a binary

```bash
make build
./bin/server
```

## API Endpoints

Base path: `/api`

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/auth/signup` | — | Register a new user |
| `POST` | `/auth/signin` | — | Sign in, receive JWT token |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/users/me` | Bearer | Get current user profile |
| `PATCH` | `/users/me` | Bearer | Update username / email / password |
| `DELETE` | `/users/me` | Bearer | Delete current user account |

### Posts

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/posts` | — | List all posts (newest first) |
| `GET` | `/posts/:id` | — | Get a single post |
| `POST` | `/posts` | Bearer | Create a post |
| `PATCH` | `/posts/:id` | Bearer | Update a post (owner only) |
| `DELETE` | `/posts/:id` | Bearer | Delete a post (owner only) |

All protected routes require the header:

```
Authorization: Bearer <token>
```

### API Documentation

Interactive Swagger UI is available at:

```
http://localhost:8080/swagger/
```

## Response Format

All responses follow a consistent envelope:

```json
{
  "success": true,
  "data": { ... }
}
```

On error:

```json
{
  "success": false,
  "error": "description of the error"
}
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make run` | Generate Swagger docs and run the server |
| `make build` | Generate Swagger docs and build binary to `bin/server` |
| `make swagger` | Regenerate Swagger documentation |
| `make tidy` | Run `go mod tidy` |
| `make install-swag` | Install the `swag` CLI tool |
| `make docker-up` | Start services with Docker Compose |
| `make docker-down` | Stop and remove Docker Compose services |
| `make docker-logs` | Tail app container logs |
| `make docker-build` | Build the Docker image |
