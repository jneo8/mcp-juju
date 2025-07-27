# MCP Juju

A Model Context Protocol (MCP) server that provides integration with [Juju](https://github.com/juju/juju).

**Note: currently only support 3.6**

## Overview

This MCP server exposes Juju functionality through MCP tools, allowing AI assistants to manage Juju controllers, models, and applications.

## Features

This MCP server supports most of the features that the Juju CLI provides (100+ commands available).

## Quick Start

### Prerequisites

- Go 1.24 or later
- Juju CLI installed and configured

### Installation

```bash
git clone https://github.com/jneo8/mcp-juju.git
cd mcp-juju
go build .
```

### Running

```bash
# Run the MCP server
make run

# Or run with custom options
./mcp-juju --port 8080 --debug
```

The server will start on `http://localhost:8080/mcp` by default.

### Configuration

Environment variables (prefixed with `MCP_JUJU_`):
- `MCP_JUJU_PORT`: Server port (default: 8080)
- `MCP_JUJU_DEBUG`: Enable debug mode (default: false)
- `MCP_JUJU_ENDPOINT`: Endpoint path (default: /mcp)

## Usage

Once running, the MCP server provides tools for all Juju CLI operations:

- `add-model`: Add a new model
- `status`: Get Juju status
- `deploy`: Deploy applications
- `add-unit`: Scale applications
- `config`: Configure applications
- `bootstrap`: Initialize a cloud environment
- `add-relation`: Create relations between applications
- And all other Juju CLI commands

## Development

### Build
```bash
go build .
```

### Test
```bash
make test
```

### Format & Lint
```bash
make lint
```

## Architecture

The project follows a clean architecture with:
- **cmd/**: CLI interface using Cobra
- **config/**: Configuration management
- **pkg/application/**: MCP server application logic
- **pkg/jujuadapter/**: Juju command adapter that converts Juju CLI commands to MCP tools
