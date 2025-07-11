# Go HTTP API - Hello World

A simple Go HTTP API server with multiple endpoints including a "Hello World" endpoint.

## Features

- **Root endpoint** (`/`) - Welcome message
- **Hello World endpoint** (`/hello`) - Returns a "Hello, World!" message
- **Health check endpoint** (`/health`) - Service health status
- JSON responses with timestamps
- Error handling

## Prerequisites

- Go 1.21 or later

## Getting Started

### 1. Initialize the module (if not already done)
```bash
go mod init hackathon-2025
```

### 2. Run the server
```bash
# Using go run
go run ./cmd/api

# Or using make
make run
```

The server will start on port 8080.

### 3. Test the endpoints

You can test the API using curl or your browser:

#### Root endpoint
```bash
curl http://localhost:8080/
```

#### Hello World endpoint
```bash
curl http://localhost:8080/hello
```

#### Health check endpoint
```bash
curl http://localhost:8080/health
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message |
| GET | `/hello` | Hello World message |
| GET | `/health` | Health check |

## Response Format

All endpoints return JSON responses with the following structure:

```json
{
  "message": "Hello, World!",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "success"
}
```

## Development

### Building the application
```bash
# Using go build
go build -o bin/api-server ./cmd/api

# Or using make
make build
```

### Running tests
```bash
# Using go test
go test ./...

# Or using make
make test
```

### Code formatting
```bash
# Using go fmt
go fmt ./...

# Or using make
make fmt
```

## Available Make Commands

The project includes a Makefile with common development tasks:

```bash
make build  # Build the application
make run    # Run the application
make test   # Run tests
make clean  # Clean build artifacts
make fmt    # Format code
make lint   # Lint code (requires golangci-lint)
make deps   # Install dependencies
make help   # Show all available commands
```

## Project Structure

```
hackathon-2025/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── pkg/
│   └── handlers/
│       └── handlers.go      # HTTP handlers
├── internal/                # Private application code
├── go.mod                  # Go module file
├── Makefile                # Build and development tasks
├── .gitignore              # Git ignore file
└── README.md               # This file
```

## Next Steps

- Add more endpoints
- Implement middleware for logging, CORS, etc.
- Add database integration
- Add authentication
- Add tests
- Add Docker support