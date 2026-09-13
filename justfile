# Task runner for mcp-juju. Run `just` to list recipes.

set shell := ["bash", "-euo", "pipefail", "-c"]

# List available recipes
[private]
default:
    @just --list --unsorted

# --- Development ---------------------------------------------------------

# Run the stdio MCP server (default transport)
run:
    go run . --server-type stdio

# Run the Streamable HTTP MCP server
run-http:
    go run . --server-type http

# Run the stdio MCP server
run-stdio: run

# Run the stdio MCP server with debug logging
run-debug:
    go run . --server-type stdio --debug

# Format Go code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Format and vet
lint: fmt vet

# Build the binary
build:
    go build .

# --- Testing -------------------------------------------------------------

# Run Go unit tests
test-unit:
    go test ./...

# Run tests and open the HTML coverage report
test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

# Regenerate testify mocks (mockery is pinned as a go.mod tool dependency)
mocks:
    go tool mockery

# Run functional tests against a bootstrapped Juju controller (needs uv), e.g. just test-functional -k test_status
[positional-arguments]
test-functional *ARGS:
    uv run --project tests/functional --group functional pytest tests/functional "$@"
