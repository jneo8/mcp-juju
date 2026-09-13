# CLI reference

`mcp-juju` starts an MCP server that exposes the Juju 3.6 CLI as tools. It takes no positional arguments.

```
mcp-juju [flags]
```

## Flags

Every flag can also be set through an environment variable: prefix `MCP_JUJU_`, flag name in upper case, dashes replaced by underscores (`--server-type` → `MCP_JUJU_SERVER_TYPE`). A flag on the command line takes precedence over the environment variable.

### General

| Flag | Type | Default | Environment variable | Description |
|---|---|---|---|---|
| `--server-type` | `string` | `stdio` | `MCP_JUJU_SERVER_TYPE` | Transport. `stdio` speaks MCP on stdin/stdout; `http` starts a Streamable HTTP server. Any other value is rejected at startup. |
| `--tool-names` | `string` list | empty (all commands) | `MCP_JUJU_TOOL_NAMES` | Comma-separated allowlist of tool names to register. Names are Juju command names as listed in [Command reference](commands.md). Unknown names fail at startup. |
| `--read-only` | `bool` | `false` | `MCP_JUJU_READ_ONLY` | Register only tools that do not modify state and enforce read-only invocation shapes; see [Read-only mode](read-only-mode.md). |
| `--debug` | `bool` | `false` | `MCP_JUJU_DEBUG` | Log at debug level. Logs always go to stderr, at info level otherwise. |

### HTTP transport

These flags apply only with `--server-type http`.

| Flag | Type | Default | Environment variable | Description |
|---|---|---|---|---|
| `--host` | `string` | `127.0.0.1` | `MCP_JUJU_HOST` | Interface to bind. A non-loopback value (including `0.0.0.0`, `::` and an empty string) requires `--auth-token` or `--allow-no-auth`. |
| `--port` | `int` | `8080` | `MCP_JUJU_PORT` | TCP port. Must be between 1 and 65535. |
| `--endpoint` | `string` | `/mcp` | `MCP_JUJU_ENDPOINT` | Path of the MCP endpoint. Other paths return `404`. |
| `--auth-token` | `string` | empty | `MCP_JUJU_AUTH_TOKEN` | Static bearer token. When set, every request must carry `Authorization: Bearer <token>`; otherwise the server responds `401` with `WWW-Authenticate: Bearer realm="mcp-juju"`. |
| `--allow-no-auth` | `bool` | `false` | `MCP_JUJU_ALLOW_NO_AUTH` | Permit a non-loopback `--host` without `--auth-token`. A warning is logged at startup. |
| `--cors-origins` | `string` list | empty (no CORS headers) | `MCP_JUJU_CORS_ORIGINS` | Comma-separated browser origins allowed by CORS. Requests from other origins receive no `Access-Control-Allow-Origin` header. |
| `--tls-cert` | `string` | empty | `MCP_JUJU_TLS_CERT` | PEM certificate file. Enables HTTPS; must be set together with `--tls-key`. |
| `--tls-key` | `string` | empty | `MCP_JUJU_TLS_KEY` | PEM private key file. Must be set together with `--tls-cert`. |

## Startup validation

The process exits with status `1` and prints the reason when:

| Condition | Message contains |
|---|---|
| `--server-type` is neither `http` nor `stdio` | `invalid server type` |
| `--server-type http` and `--port` is outside 1–65535 | `invalid port` |
| `--server-type http`, non-loopback `--host`, no `--auth-token`, no `--allow-no-auth` | `refusing to serve HTTP on non-loopback host` |
| Exactly one of `--tls-cert` / `--tls-key` is set | `--tls-cert and --tls-key must be set together` |
| A name in `--tool-names` is not a known command | `unknown command` |

## Runtime behaviour

| Aspect | Behaviour |
|---|---|
| Juju client store | The server reads the same client store as the `juju` CLI (`$JUJU_DATA`, default `~/.local/share/juju`). Controllers, models and credentials must already be registered there. |
| Current model | There is no `switch` tool. Every model-scoped tool accepts a `model` argument; without it, the client store's current model is used. |
| Standard input | Juju commands run with an empty stdin. Commands that would prompt fail instead; pass their `no-prompt`, `yes` or `force` flag. |
| Logging | zerolog JSON lines on stderr. Never on stdout. |
| Shutdown (HTTP) | `SIGINT` and `SIGTERM` trigger a graceful shutdown with a 10 second timeout. |
| Request header timeout (HTTP) | 10 seconds. |
| DNS rebinding protection (HTTP) | Requests arriving over a loopback connection whose `Host` header is not a localhost value receive `403`. Always on. |

## Exit status

| Status | Meaning |
|---|---|
| `0` | Server stopped normally (stdio peer closed the stream, or HTTP server shut down on signal). |
| `1` | Startup validation failed, or the server returned an error. |
