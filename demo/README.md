# How to run the mcp-juju demo

The demo stands up a small Juju deployment through the MCP server itself, checks the server against it, and tears it down again. Three scripts, one per phase, so the deployment can stay around while you show the server to an MCP client.

The scripts follow the same pattern as the functional tests: the MCP server does every change through its tools, and [Jubilant](https://documentation.ubuntu.com/jubilant/) observes the result through the `juju` CLI.

## Prerequisites

- [uv](https://docs.astral.sh/uv/) and Go 1.26 (the scripts build `mcp-juju` from source unless `MCP_JUJU_BINARY` or `--binary` points at a binary).
- A bootstrapped Juju controller on a machine cloud, for example `juju bootstrap localhost lxd`. The default charm, `postgresql`, is a machine charm.

## 1. Deploy

```bash
just demo-deploy
```

Creates the model `mcp-juju-demo` with the MCP `add-model` tool, deploys `postgresql` as `db` with the MCP `deploy` tool, and waits until the unit is active. Rerunning reuses whatever already exists. When it finishes it prints how to point an MCP client at the same binary, with and without `--read-only`.

Watch progress from another terminal with:

```bash
juju status --watch 2s -m mcp-juju-demo
```

## 2. Verify

```bash
just demo-verify
```

Starts the server three ways and, for each check, prints the name, an excerpt of what the server returned (status, models, config, rejection messages, dry-run output, HTTP headers), and the verdict. Pass `--full` to print tool output in full instead of excerpts.

- **stdio, every tool**: instructions, tool list and annotations, the `juju://status-doc` resource, and `status`, `models`, `show-application` and `config` results compared with `juju` output.
- **stdio, `--read-only`**: mutating tools are absent, queries work, a `config` write and `status --output` are rejected before reaching Juju, and `remove-unit` is accepted only with `--dry-run`.
- **Streamable HTTP with a bearer token**: a request without the token gets 401, and the tool list and `status` match the stdio server.

Only read-only operations run, so it can be repeated at any time. The exit status is non-zero when a check fails.

## 3. Clean

```bash
just demo-clean                  # destroy the model and its storage
just demo-clean --remove-build   # also delete demo/.build/mcp-juju
```

Destroys the model with the MCP `destroy-model` tool and waits until Juju no longer lists it.

## Options

Every script accepts `--model`, `--controller` and `--binary`; `deploy` also accepts `--charm`. The same values can be set with `MCP_JUJU_DEMO_MODEL`, `MCP_JUJU_DEMO_CONTROLLER`, `MCP_JUJU_BINARY` and `MCP_JUJU_DEMO_CHARM`; `MCP_JUJU_DEMO_APP` changes the application name. Pass script options after the recipe name:

```bash
just demo-deploy --controller lxd --charm postgresql
just demo-verify --controller lxd
just demo-clean --controller lxd
```
