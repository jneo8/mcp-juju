# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server implementation in Go that provides integration with Juju (Ubuntu's application modeling tool). The server exposes Juju functionality through MCP tools, allowing AI assistants to manage Juju controllers, models, and applications through 100+ Juju CLI commands.

**Important**: Currently only supports Juju 3.6.

## Architecture

The codebase follows a clean architecture pattern with clear separation of concerns:

### Core Structure
- **main.go**: Entry point that delegates to the cmd package
- **cmd/**: Command-line interface implementation using Cobra
- **config/**: Configuration management with environment variable support
- **pkg/application/**: MCP server application logic and HTTP server setup
- **pkg/jujuadapter/**: Juju command adapter that converts all Juju CLI commands to MCP tools

### Key Components

#### Application Layer (`pkg/application/`)
- **app.go**: Main application structure that initializes the MCP server and registers tools
- **server.go**: HTTP server setup for MCP communication using Server-Sent Events (SSE)

#### Juju Adapter Layer (`pkg/jujuadapter/`)
- **adapter.go**: Core adapter that converts Juju commands to MCP tools with automatic flag detection
- **command.go**: Command interface and implementation with comprehensive command list (232 commands)
- **factory.go**: Command factory for creating Juju command instances

### MCP Tool Generation System

The adapter automatically generates MCP tools from Juju commands by:
1. Reflecting on command flag sets to determine parameter types (boolean, number, string)
2. Building enhanced descriptions from Juju command documentation
3. Supporting both positional arguments and named flags
4. Converting flag types to appropriate MCP tool options

## Common Development Commands

### Build and Run
```bash
# Build the application
go build .

# Run the application (uses Makefile)
make run

# Run with debug logging
make run-debug

# Run with custom port
./mcp-juju --port 8080 --debug
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./pkg/jujuadapter/
```

### Code Quality
```bash
# Format and lint code
make lint

# Format code only
make fmt

# Run go vet only
make vet
```

## Environment Variables

The application supports configuration through environment variables with the prefix `MCP_JUJU_`:

- `MCP_JUJU_PORT`: Server port (default: 8080)
- `MCP_JUJU_DEBUG`: Enable debug mode (default: false)
- `MCP_JUJU_ENDPOINT`: Endpoint path (default: /mcp)

## Key Dependencies

- **github.com/juju/juju**: Core Juju client library (custom fork with 3.6 support)
- **github.com/mark3labs/mcp-go**: MCP protocol implementation
- **github.com/spf13/cobra**: CLI framework
- **github.com/spf13/viper**: Configuration management with automatic env binding
- **github.com/rs/zerolog**: Structured logging

## Development Notes

### Command Registration System
- All 232 Juju commands are registered in `pkg/jujuadapter/command.go` in the exact order as the official Juju CLI
- Commands are dynamically converted to MCP tools using reflection on their flag sets
- The adapter supports positional arguments through a special "args" array parameter

### MCP Communication
- Uses Server-Sent Events (SSE) for MCP communication over HTTP
- Server runs on configurable port with `/mcp` endpoint by default
- All Juju operations are context-aware and support cancellation

### Flag Type Detection
The adapter uses reflection to detect Juju command flag types and converts them to appropriate MCP tool parameters:
- `boolValue` → boolean parameter
- `intValue`/`float64Value` → number parameter  
- String and unknown types → string parameter