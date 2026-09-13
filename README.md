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
just run

# Run the Streamable HTTP server
./mcp-juju --server-type http --port 8080

# Expose only a subset of commands, with debug logging on stderr
./mcp-juju --tool-names status,deploy,config --debug
```

In HTTP mode the server listens on `http://localhost:8080/mcp` by default.

### Configuration

Environment variables (prefixed with `MCP_JUJU_`):
- `MCP_JUJU_SERVER_TYPE`: `stdio` (default) or `http`
- `MCP_JUJU_TOOL_NAMES`: Comma-separated list of command IDs to expose (default: all)
- `MCP_JUJU_DEBUG`: Enable debug logging on stderr (default: false)
- `MCP_JUJU_HOST`: Interface for HTTP mode (default: 127.0.0.1)
- `MCP_JUJU_PORT`: Port for HTTP mode (default: 8080)
- `MCP_JUJU_ENDPOINT`: Endpoint path for HTTP mode (default: /mcp)
- `MCP_JUJU_AUTH_TOKEN`: Bearer token HTTP clients must send
- `MCP_JUJU_ALLOW_NO_AUTH`: Allow a non-loopback host without a token (default: false)
- `MCP_JUJU_CORS_ORIGINS`: Comma-separated browser origins allowed by CORS
- `MCP_JUJU_TLS_CERT`, `MCP_JUJU_TLS_KEY`: Enable HTTPS

### HTTP mode security

The HTTP server binds to `127.0.0.1` by default. Binding another host (`--host 0.0.0.0`) is refused unless `--auth-token` is set, or `--allow-no-auth` is passed explicitly, because every Juju command would otherwise be reachable unauthenticated. With a token, every request must carry `Authorization: Bearer <token>`. Requests over loopback whose `Host` header is not a localhost value are rejected (DNS rebinding protection), and CORS headers are only emitted for origins listed in `--cors-origins`.

```bash
mcp-juju --server-type http --host 0.0.0.0 --auth-token "$(openssl rand -hex 32)" \
  --cors-origins https://app.example.com --tls-cert cert.pem --tls-key key.pem
```

## Usage

Once running, the MCP server provides tools for all Juju CLI operations. Each tool shares the name of its Juju command and follows the same conventions:

- Positional arguments go in the `args` array; every other property is a long flag name.
- Commands with `--format` default to `json`, and the parsed object is also returned as MCP `structuredContent`.
- Tools carry MCP annotations (read-only, destructive, idempotent, open-world) so clients can ask for confirmation before destructive commands.
- Failed commands return an `isError` result containing the Juju error text.
- Commands never read stdin, so confirmation prompts must be skipped with flags such as `no-prompt`.
- The full help text of each command is available as the resource `juju://<tool>-doc`.

Examples:

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
just test-unit
```

### Format & Lint
```bash
just lint
```

### Functional tests

The functional tests in `tests/functional/` start the real `mcp-juju` binary, drive it through the [MCP Python SDK](https://github.com/modelcontextprotocol/python-sdk), and verify the results with [Jubilant](https://documentation.ubuntu.com/jubilant/). One temporary model is created and the charm is deployed once through the MCP `deploy` tool; every test then verifies a tool or resource against it. They need [uv](https://docs.astral.sh/uv/) and a bootstrapped Juju controller; [pytest-jubilant](https://github.com/canonical/pytest-jubilant) creates a temporary model per test module and destroys it afterwards.

```bash
juju bootstrap localhost lxd          # once; any machine cloud works with the default charm
just test-functional                       # run everything
just test-functional -k test_status        # run a subset
just test-functional --no-juju-teardown    # keep the model for inspection
just test-functional --juju-controller lxd --test-charm postgresql
```

Set `MCP_JUJU_BINARY` to test a prebuilt binary instead of building from source.

## Architecture

The project follows a clean architecture with:
- **cmd/**: CLI interface using Cobra
- **config/**: Configuration management
- **pkg/application/**: MCP server application logic
- **pkg/jujuadapter/**: Juju command adapter that converts Juju CLI commands to MCP tools
