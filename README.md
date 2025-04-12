# Backend Dragonhak

A robust Go backend service for managing craftsmen, workshops, and user interactions.

## Features

- **Authentication & Authorization**

  - JWT-based authentication
  - Role-based access control
  - Rate limiting for auth endpoints
  - Secure password handling with bcrypt

- **User Management**

  - User registration and profile management
  - Email verification
  - Password reset functionality
  - User roles (Admin, Craftsman, User)

- **Craftsman Features**

  - Craftsman profile management
  - Workshop creation and management
  - Specialties and skills tracking
  - Rating and review system

- **Workshop Management**

  - Workshop scheduling
  - Participant management
  - Materials and requirements tracking
  - Location and contact information

- **Security**
  - Rate limiting with Redis
  - Secure password policies
  - Input validation
  - CORS protection

## Prerequisites

- Go 1.23 or later
- MongoDB
- Redis
- Make (for development)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/backend-dragonhak.git
cd backend-dragonhak
```

2. Install dependencies:

```bash
make deps
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

## Development

### Available Make Commands

- `make` - Run deps, lint, test, and build
- `make build` - Build the application
- `make clean` - Clean build files
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage report
- `make run` - Build and run the application
- `make deps` - Install dependencies
- `make lint` - Run linter
- `make build-linux` - Build for Linux
- `make help` - Show help message

### Code Quality

The project uses golangci-lint for code quality checks. The linter configuration is in `.golangci.yml`.

### Testing

Run tests with:

```bash
make test
```

For coverage report:

```bash
make test-coverage
```

## API Documentation

API documentation is available at `/docs` when running the server.

## Environment Variables

Required environment variables:

- `MONGODB_URI` - MongoDB connection string
- `REDIS_ADDR` - Redis server address
- `JWT_SECRET` - JWT signing secret
- `RATE_LIMIT_WINDOW` - Rate limit window in seconds
- `RATE_LIMIT_MAX_REQUESTS` - Maximum requests per window

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
