# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server implementation in Go that provides integration with Juju (Ubuntu's application modeling tool). The server exposes Juju functionality through MCP tools, allowing AI assistants to manage Juju controllers, models, and applications.

## Architecture

The codebase follows a clean architecture pattern with the following key components:

### Core Structure
- **main.go**: Entry point that delegates to the cmd package
- **cmd/**: Command-line interface implementation using Cobra
- **config/**: Configuration management with environment variable support
- **pkg/application/**: MCP server application logic and tool handlers
- **pkg/jujuclient/**: Juju API client wrapper

### Key Components

#### Application Layer (`pkg/application/`)
- **app.go**: Main application structure that initializes the MCP server and registers tools
- **handlers.go**: MCP tool handlers for Juju operations
- **arguments.go**: Tool argument definitions
- **sse.go**: Server-Sent Events implementation for MCP communication

#### Juju Client Layer (`pkg/jujuclient/`)
- **client.go**: Main Juju client implementation
- **types.go**: Data structures for Juju responses
- **logger.go**: Logging adapter for Juju client
- **errors.go**: Error handling utilities

### MCP Tools Available
1. **listControllers**: List all Juju controllers
2. **listModels**: List models for a controller
3. **getStatus**: Get Juju status for a model
4. **getApplicationConfig**: Get application configuration
5. **setApplicationConfig**: Set application configuration

## Common Development Commands

### Build and Run
```bash
# Build the application
go build -o mcp-juju

# Run the application
go run main.go

# Run with custom host/port
go run main.go --host localhost --port 8080

# Run with debug logging
go run main.go --debug
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/jujuclient/
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (if available)
go vet ./...

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Environment Variables

The application supports configuration through environment variables with the prefix `MCP_JUJU_`:

- `MCP_JUJU_HOST`: Server host (default: localhost)
- `MCP_JUJU_PORT`: Server port (default: 8080)
- `MCP_JUJU_DEBUG`: Enable debug mode (default: false)

## Dependencies

Key dependencies include:
- **github.com/juju/juju**: Core Juju client library
- **github.com/mark3labs/mcp-go**: MCP protocol implementation
- **github.com/spf13/cobra**: CLI framework
- **github.com/spf13/viper**: Configuration management
- **github.com/rs/zerolog**: Structured logging

## Development Notes

- The application uses Server-Sent Events (SSE) for MCP communication
- All Juju operations are context-aware and support cancellation
- Error handling follows Juju's error package patterns
- Configuration supports both CLI flags and environment variables
- The client automatically handles controller/model defaults when not specified