# MCP Juju

A Model Context Protocol (MCP) server that provides integration with [Juju](https://github.com/juju/juju).

**Note: currently only support 3.6**

## Overview

This MCP server exposes Juju functionality through MCP tools, allowing AI assistants to manage Juju controllers, models, and applications.

## Features

This MCP server supports most of the features that the Juju CLI provides (161 commands available).

## Quick Start

### Prerequisites

- Go 1.26 or later
- Juju CLI installed and configured

### Installation

#### From Source

```bash
git clone https://github.com/jneo8/mcp-juju.git
cd mcp-juju
go build .
```

#### From Snap Package

**Note**: The snap package includes a daemon service that requires experimental user daemon support.

1. Enable experimental user daemons:
```bash
sudo snap set system experimental.user-daemons=true
```

2. Install the snap:
```bash
sudo snap install mcp-juju
# Or install from local build:
sudo snap install ./mcp-juju_*.snap --dangerous
```

3. Connect the required interface for Juju data access:
```bash
sudo snap connect mcp-juju:dot-local-share-juju
```

4. Start the daemon (optional):
```bash
sudo snap start mcp-juju.mcp-juju-daemon
```

### Running

```bash
# Run the MCP server over stdio (default)
make run

# Run the Streamable HTTP server
./mcp-juju --server-type http --port 8080

# Expose only a subset of commands, with debug logging on stderr
./mcp-juju --tool-names status,deploy,config --debug
```

In HTTP mode the server listens on `http://localhost:8080/mcp` by default.

### Configuration

Environment variables (prefixed with `MCP_JUJU_`):
- `MCP_JUJU_PORT`: Server port for HTTP mode (default: 8080)
- `MCP_JUJU_DEBUG`: Enable debug logging (default: false)
- `MCP_JUJU_ENDPOINT`: Endpoint path for HTTP mode (default: /mcp)
- `MCP_JUJU_SERVER_TYPE`: `stdio` (default) or `http`
- `MCP_JUJU_TOOL_NAMES`: Comma-separated list of command IDs to expose (default: all)

## Usage

Once running, the MCP server provides tools for all Juju CLI operations:

- `add-model`: Add a new model
- `status`: Get Juju status
- `deploy`: Deploy applications
- `add-unit`: Scale applications
- `config`: Configure applications
- `bootstrap`: Initialize a cloud environment
- `integrate`: Create relations between applications
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
