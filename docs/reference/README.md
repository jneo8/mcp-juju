# Reference

Precise descriptions of what `mcp-juju` exposes and accepts. For task-oriented instructions see the main [README](../../README.md).

| Page | Covers |
|---|---|
| [CLI reference](cli.md) | Flags, environment variables, startup validation, runtime behaviour, exit status |
| [MCP interface reference](mcp-interface.md) | Server identity, tool definitions and schemas, argument handling, results, resources, resource templates, transports |
| [Command reference](commands.md) | All 151 tools with their annotations and read-only availability; commands that are not exposed |
| [Read-only mode reference](read-only-mode.md) | Registration rules, host-write flags, per-command policies, rejection messages |

`commands.md` is generated from the adapter's own tables (`command_defs.go`, `annotations.go`, `readonly.go`) by `just docs-commands`; a unit test fails when it is out of date.
