# Task runner for mcp-juju. Run `just` to list recipes by group.

set shell := ["bash", "-euo", "pipefail", "-c"]

# List available recipes
[private]
default:
    @just --list --unsorted

# Build the binary
[group('dev')]
build:
    go build .

# Run the stdio MCP server (default transport)
[group('dev')]
run:
    go run . --server-type stdio

# Run the stdio MCP server
[group('dev')]
run-stdio: run

# Run the stdio MCP server with debug logging
[group('dev')]
run-debug:
    go run . --server-type stdio --debug

# Run the Streamable HTTP MCP server
[group('dev')]
run-http:
    go run . --server-type http

# Open the MCP Inspector web UI on the stdio server (needs npx), e.g. just inspect --read-only
[group('dev')]
[positional-arguments]
inspect *ARGS: build
    #!/usr/bin/env bash
    set -euo pipefail
    just --justfile "{{ justfile() }}" _inspect-examples "$@"
    exec npx @modelcontextprotocol/inspector ./mcp-juju -- --server-type stdio "$@"

# Call one tool through the Inspector CLI (uses ./mcp-juju as built; run `just build` to refresh it)
[group('dev')]
[positional-arguments]
call TOOL ARGS_JSON='{}' *FLAGS:
    #!/usr/bin/env bash
    set -euo pipefail
    [ -x ./mcp-juju ] || go build .
    tool="$1"; args_json="$2"; shift 2
    exec npx @modelcontextprotocol/inspector --cli ./mcp-juju --server-type stdio "$@" -- \
        --method tools/call --tool-name "$tool" --tool-args-json "$args_json"

# Print read-only example calls for the demo model, as JSON for the web UI and as `just call` lines
[private]
[positional-arguments]
_inspect-examples *FLAGS:
    #!/usr/bin/env bash
    set -euo pipefail
    model="${MCP_JUJU_DEMO_MODEL:-mcp-juju-demo}"
    app="${MCP_JUJU_DEMO_APP:-db}"
    accepted=(
      "status|{\"model\": \"$model\"}"
      "models|{}"
      "show-model|{\"args\": [\"$model\"]}"
      "show-application|{\"model\": \"$model\", \"args\": [\"$app\"]}"
      "show-unit|{\"model\": \"$model\", \"args\": [\"$app/0\"]}"
      "config|{\"model\": \"$model\", \"args\": [\"$app\"]}"
      "config|{\"model\": \"$model\", \"args\": [\"$app\", \"profile\"]}"
      "model-config|{\"model\": \"$model\"}"
      "controller-config|{}"
      "machines|{\"model\": \"$model\"}"
      "storage|{\"model\": \"$model\"}"
      "remove-unit|{\"model\": \"$model\", \"args\": [\"$app/0\"], \"dry-run\": true}"
    )
    rejected=(
      "config|{\"model\": \"$model\", \"args\": [\"$app\", \"profile=testing\"]}"
      "status|{\"model\": \"$model\", \"output\": \"/tmp/status.json\"}"
      "remove-unit|{\"model\": \"$model\", \"args\": [\"$app/0\"]}"
    )
    echo "Example calls for the demo model (set MCP_JUJU_DEMO_MODEL / MCP_JUJU_DEMO_APP to change the values)."
    echo
    echo "Web UI: Tools tab -> select the tool -> 'Edit as JSON' -> paste the JSON -> Execute Tool"
    for e in "${accepted[@]}"; do printf '  %-18s %s\n' "${e%%|*}" "${e#*|}"; done
    echo "  Rejected with --read-only (isError, never reaches Juju):"
    for e in "${rejected[@]}"; do printf '  %-18s %s\n' "${e%%|*}" "${e#*|}"; done
    echo "  Resources tab: status-doc, config-doc, or any <tool>-doc"
    echo
    echo "Terminal (Inspector CLI mode, same flags as this run: $*):"
    for e in "${accepted[@]}" "${rejected[@]}"; do printf "  just call %s '%s' %s\n" "${e%%|*}" "${e#*|}" "$*"; done
    echo

# Format Go code
[group('quality')]
fmt:
    go fmt ./...

# Run go vet
[group('quality')]
vet:
    go vet ./...

# Format and vet
[group('quality')]
lint: fmt vet

# Run Go unit tests
[group('test')]
test-unit:
    go test ./...

# Run tests and open the HTML coverage report
[group('test')]
test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

# Run functional tests against a bootstrapped Juju controller (needs uv), e.g. just test-functional -k test_status
[group('test')]
[positional-arguments]
test-functional *ARGS:
    @echo "Tip: watch the temporary models from another terminal with:"
    @echo "  watch -c \"juju models | grep jubilant | awk '{print \\\$1}' | tr -d '*' | xargs -I {} juju status --color -m {}\""
    @echo
    uv run --project tests/functional --group functional pytest tests/functional "$@"

# Regenerate testify mocks (mockery is pinned as a go.mod tool dependency)
[group('test')]
mocks:
    go tool mockery

# Deploy the demo model and charm through the MCP server (needs uv and a bootstrapped controller)
[group('demo')]
[positional-arguments]
demo-deploy *ARGS:
    uv run --project demo python demo/deploy.py "$@"

# Verify the server against the demo deployment: stdio, --read-only and HTTP
[group('demo')]
[positional-arguments]
demo-verify *ARGS:
    uv run --project demo python demo/verify.py "$@"

# Destroy the demo model through the MCP server
[group('demo')]
[positional-arguments]
demo-clean *ARGS:
    uv run --project demo python demo/clean.py "$@"

# Regenerate docs/reference/commands.md from the adapter's command tables
[group('docs')]
docs-commands:
    go run ./tools/gendocs
