# MCP interface reference

The server implements the Model Context Protocol through `github.com/mark3labs/mcp-go` v1.0.0. It speaks protocol version `2026-07-28` and negotiates down to clients that use the `initialize` handshake (`2025-11-25`, `2025-06-18`, `2025-03-26`, `2024-11-05`).

## Server identity

| Field | Value |
|---|---|
| `serverInfo.name` | `mcp-juju-server` |
| `serverInfo.version` | Build version (`config.Version`) |
| `capabilities.tools.listChanged` | `true` (advertised by the library; the tool list never changes at runtime) |
| `capabilities.resources` | present, without `subscribe` or `listChanged` |
| `capabilities.logging` | not advertised |
| `instructions` | Usage guidance for the model; extended with a read-only paragraph when `--read-only` is set |

## Tools

One tool per Juju command; the tool name is the Juju command name (`status`, `deploy`, `integrate`, ...). See [Command reference](commands.md) for the complete list.

### Tool definition

| Field | Content |
|---|---|
| `name` | Juju command name |
| `title` | `juju <name>` |
| `description` | Juju purpose text, the Juju usage line prefixed with `Arguments:` when the command takes positional arguments, and `Full help: read the resource juju://<name>-doc.` In read-only mode, tools with a policy carry an extra `Read-only mode: ...` sentence. |
| `annotations` | `title`, `readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`, all explicitly set; values per [Command reference](commands.md) |
| `inputSchema` | JSON Schema object built from the command's flag set, see below |

### Input schema

| Property | Type | Meaning |
|---|---|---|
| `args` | `array` of `string` | Positional arguments, in CLI order. Empty strings are dropped. |
| `<flag>` | see table below | One property per long flag. Short aliases that share a value with a long flag (`-B` / `--base`) are folded into the long name. |

Flag types map to schema types by the flag's Go value type:

| Juju flag value | Schema | Notes |
|---|---|---|
| Any value implementing gnuflag `IsBoolFlag()` (`boolValue`, `AutoBoolValue`, `optBoolValue`) | `boolean` | `default` set only when the Juju default is `true` or `false` |
| `intValue`, `int64Value`, `uintValue`, `uint64Value`, `float64Value` | `number` | `default` set when the Juju default parses as a number |
| `durationValue` | `string` | Description ends with `(duration such as 30s, 5m or 1h)` |
| `cmd.FileVar` | `string` | Description ends with `(path to a local file)` |
| `cmd.StringsValue` | `array` of `string` | Joined with `,` and passed once |
| `cmd.AppendStringsValue` | `array` of `string` | Passed once per element |
| `formatterValue` (`--format`) | `string` with `enum` | Choices parsed from the usage text; `default` is `json` when available, otherwise the Juju default |
| Anything else | `string` | `default` set when the Juju default is non-empty |

Flag descriptions are the Juju usage strings.

### Argument handling on `tools/call`

| Input | Handling |
|---|---|
| Property not in the schema | Ignored |
| Empty string, `null`, empty array | Treated as unset |
| `boolean`, `number` | Converted to text and passed to the flag |
| `array` on a non-array flag | Error result |
| `format` omitted on a command that supports `json` | `--format json` is applied |

### Result

| Case | `isError` | `content` | `structuredContent` |
|---|---|---|---|
| Success, `--format json` in effect, stdout parses as a JSON object | `false` | `[0]`: text, exactly the JSON stdout. `[1]` (only when stderr is non-empty): text, the stderr output | The parsed object |
| Success, otherwise | `false` | One text block: stdout, then stderr, separated by a newline | absent |
| Juju command failed, flag rejected, argument invalid, read-only rejection | `true` | One text block: `juju <name> failed: <error>`, then a blank line and whatever the command printed | absent |
| Tool name unknown | — | JSON-RPC error (`-32602`) | — |

Command execution errors are never JSON-RPC errors.

### Tool list order

`tools/list` returns tools in Juju CLI registration order, unchanged between calls. `--tool-names` preserves the order given on the command line.

## Resources

### Documentation resources

One resource per registered tool.

| Field | Value |
|---|---|
| `uri` | `juju://<name>-doc` |
| `name` | `<name>-doc` |
| `description` | `Documentation for <name> command` |
| `mimeType` | `text/markdown` |

Content is Markdown: `# <name>`, the purpose, an `## Arguments` section when the command takes positional arguments, and a `## Details` section with the full Juju help text.

Reading a URI that does not match a registered resource returns a JSON-RPC error.

## Resource templates

### `juju://config/{application}{/config_name*}`

| Field | Value |
|---|---|
| `name` | `Juju Application Configuration` |
| `description` | `Get configuration for any Juju application or specific config key in JSON format` |
| `mimeType` | `application/json` |
| Underlying command | `config --format json <application> [<config_name>]` |
| Model | The client store's current model; the template has no model variable |
| Read-only mode | Allowed (query shape) |

Expansion examples:

| URI | Command |
|---|---|
| `juju://config/postgresql` | `config --format json postgresql` |
| `juju://config/postgresql/profile` | `config --format json postgresql profile` |

The returned `contents[0].text` is the command's stdout; `contents[0].mimeType` is `application/json` even when Juju prints an error message.

## Transport specifics

| Transport | Endpoint | Notes |
|---|---|---|
| stdio | stdin / stdout | One JSON-RPC message per line. Nothing else is written to stdout. |
| Streamable HTTP | `http[s]://<host>:<port><endpoint>` | `POST` for requests; `Accept` must include `application/json` and `text/event-stream`. Bearer authentication, CORS, TLS and DNS rebinding protection as described in [CLI reference](cli.md). |
