# Task Management System

A robust task management system built with Go, featuring user authentication, role-based access control, task management, real-time updates via WebSocket, and comprehensive chat functionality.

## Features

- **User Management**

  - User registration and authentication
  - Role-based access control (Employee/Employer)
  - JWT-based authentication
  - Secure password hashing
  - Casbin for authorization
- **Task Management**

  - Create, read, update, and delete tasks
  - Task status tracking (Pending, In Progress, Completed)
  - Task assignment to employees
  - Task filtering and sorting
  - Task summary by employee
  - Real-time task updates via WebSocket
- **Chat**

  - Real-time messaging between users
  - Direct and group chat rooms
  - Message history and pagination
  - Message read receipts
  - Message pinning/unpinning
  - Room archiving/unarchiving
  - Room muting/unmuting
  - Room joining/leaving
  - Room updates
- **API Features**

  - RESTful API design
  - Input validation using validator
  - Error handling
  - Swagger documentation
  - WebSocket support for real-time updates and chat

## Tech Stack

- **Language**: Go 1.23
- **Framework**: 
  - Chi (Router)
  - Wire (Dependency Injection)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication & Authorization**:
  - JWT
  - Casbin
- **Configuration**: Viper
- **Testing**:
  - Testify/Suite for test organization
  - Uber's gomock for mocking
  - Table-driven tests
- **Documentation**: Swagger
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── api/                 # Application entry point
├── config/                  # Configuration files
├── docs/                    # Documentation
├── internal/
│   ├── delivery/           # Delivery layer (HTTP handlers)
│   │   └── rest/
│   │       ├── handler/    # HTTP handlers
│   │       └── dtos/       # Data Transfer Objects
│   ├── domain/            # Domain layer (business logic)
│   │   ├── task/         # Task domain
│   │   └── user/         # User domain
│   ├── repositories/      # Repository layer
│   │   └── postgres/     # PostgreSQL implementation
│   ├── usecase/          # Use case layer
│   ├── mocks/            # Generated mocks for testing
│   └── server/           # Server configuration and setup
├── migrations/            # Database migrations
├── pkg/                   # Shared packages
│   ├── apperrors/        # Application errors
│   ├── cache/           # Caching utilities
│   └── utils/           # Utility functions
├── postgres-data/        # PostgreSQL data directory
├── .github/              # GitHub configuration
├── .vscode/              # VS Code configuration
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose configuration
├── go.mod                # Go module definition
└── go.sum                # Go dependencies checksum
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Docker & Docker Compose (optional)
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/personal/task-management.git
cd task-management
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:

```bash
# Using Go directly
go run cmd/api/main.go

# Using Docker
docker-compose up
```

### Running Tests

```bash
# Generate mocks
cd internal/mocks && go generate ./...

# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific test package
go test ./internal/usecase/... -v
```

## API Documentation

The API documentation is available via Swagger UI when running the application:

```
http://localhost:8080/swagger/index.html
```

### Main Endpoints

- **Authentication**

  - `POST /auth/register` - Register a new user
  - `POST /auth/login` - Login user
- **Users**

  - `GET /users` - List users
  - `GET /users/{id}` - Get user by ID
  - `PUT /users/{id}` - Update user
  - `DELETE /users/{id}` - Delete user
- **Tasks**

  - `POST /tasks` - Create task
  - `GET /tasks` - List tasks
  - `GET /tasks/{id}` - Get task by ID
  - `PUT /tasks/{id}` - Update task
  - `DELETE /tasks/{id}` - Delete task
  - `GET /tasks/employee/{id}` - Get employee tasks
  - `GET /tasks/summary` - Get task summary by employee
- **Chat**

  - **Room Management**
    - `POST /chat/rooms/direct` - Create direct chat room
    - `POST /chat/rooms/group` - Create group chat room
    - `GET /chat/rooms` - List all rooms
    - `GET /chat/rooms/{roomId}` - Get room history
    - `POST /chat/rooms/{roomId}/join` - Join a room
    - `POST /chat/rooms/{roomId}/leave` - Leave a room
    - `PUT /chat/rooms/{roomId}` - Update room details
  
  - **Message Management**
    - `GET /chat/rooms/{roomId}/messages` - Get room messages
    - `POST /chat/rooms/{roomId}/messages` - Send message
    - `POST /chat/rooms/{roomId}/messages/{messageId}/read` - Mark message as read
    - `POST /chat/rooms/{roomId}/messages/{messageId}/pin` - Pin message
    - `DELETE /chat/rooms/{roomId}/messages/{messageId}/pin` - Unpin message
  
  - **Room Actions**
    - `POST /chat/rooms/{roomId}/archive` - Archive room
    - `POST /chat/rooms/{roomId}/unarchive` - Unarchive room
    - `POST /chat/rooms/{roomId}/mute` - Mute room notifications
    - `POST /chat/rooms/{roomId}/unmute` - Unmute room notifications

- **WebSocket**

  - `WS /ws` - WebSocket connection for real-time updates and chat

## Testing

The project uses a comprehensive testing approach:

- **Unit Tests**: Using testify/suite for organized test structure
- **Mocking**: Using Uber's gomock for dependency mocking
- **Integration Tests**: End-to-end API testing
- **Table-driven Tests**: For testing multiple scenarios

### Running Tests

```bash
# Generate mocks
cd internal/mocks && go generate ./...

# Run all tests
go test ./... -v

# Run specific test package
go test ./internal/usecase/... -v
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
