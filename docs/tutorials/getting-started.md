# Tutorial: Your first Juju deployment through mcp-juju

In this tutorial you build the `mcp-juju` server, use it to deploy PostgreSQL into a Juju model, talk to the server from a few lines of Python, see what read-only mode refuses, and finally browse the server in the MCP Inspector. By the end you will know how an MCP client sees Juju: tools, arguments, results and documentation.

It takes about 30 minutes, most of it waiting for the LXD machine to come up.

## Prerequisites

- Ubuntu 22.04 or later with [LXD](https://documentation.ubuntu.com/lxd/) initialised (`lxd init --auto` is enough)
- Juju 3.6 CLI: `sudo snap install juju --channel 3.6/stable`
- Go 1.26 or later, [just](https://just.systems/), [uv](https://docs.astral.sh/uv/) and git
- Node.js 22 or later with `npx`, for the [MCP Inspector](https://github.com/modelcontextprotocol/inspector) in the last step
- A clone of this repository:

  ```bash
  git clone https://github.com/jneo8/mcp-juju.git
  cd mcp-juju
  ```

Run every command below from the repository root unless a step says otherwise.

## Step 1: Bootstrap a Juju controller

The server is a client of Juju, so it needs a controller to talk to. Bootstrap one on your local LXD:

```bash
juju bootstrap localhost lxd
```

This takes a few minutes. It ends with:

```
Bootstrap complete, controller "lxd" is now available
Controller machines are in the "controller" model
```

Confirm the CLI can see it:

```bash
juju controllers
```

```
Controller  Model    User   Access     Cloud/Region         Models  Nodes    HA  Version
lxd*        default  admin  superuser  localhost/localhost       2      1  none  3.6.x
```

## Step 2: Build the server

```bash
just build
```

This produces the `mcp-juju` binary in the repository root. Check it runs:

```bash
./mcp-juju --help
```

You should see the usage text starting with `MCP Juju` and a list of flags such as `--read-only`, `--server-type` and `--tool-names`.

## Step 3: Deploy PostgreSQL through the server

Instead of running `juju deploy` yourself, you let the MCP server do it. The demo script starts the server, calls its `add-model` tool to create the model `mcp-juju-demo`, calls its `deploy` tool to deploy the `postgresql` charm as `db`, and waits until the unit is active:

```bash
just demo-deploy
```

The first run builds a copy of the binary for the demo and then prints each tool call as it happens:

```
=== Model mcp-juju-demo ===
[   1.2s] mcp: add-model {'args': ['mcp-juju-demo'], 'no-switch': True}

=== Application db (postgresql) ===
[   4.8s] mcp: deploy {'model': 'mcp-juju-demo', 'args': ['postgresql', 'db']}
[   6.1s] waiting up to 20 min for db to become active
```

While you wait, open a second terminal and watch Juju bring the machine up:

```bash
juju status --watch 2s -m mcp-juju-demo
```

When the unit reaches `active`, the script finishes with:

```
=== Ready ===
model:   mcp-juju-demo
charm:   postgresql (rev 553)
unit:    db/0 on machine 0, active
```

followed by the commands for the later steps. Everything Juju knows about this model was created by MCP tool calls; you never ran `juju add-model` or `juju deploy`.

## Step 4: Call a tool yourself

Now talk to the server directly. The demo directory ships a small synchronous MCP client, so you can drive the server from a few lines of Python. Change into the demo directory:

```bash
cd demo
```

and run this script. It starts the server over stdio, counts its tools, and calls the `status` tool the same way an AI assistant would:

```bash
uv run python - <<'EOF'
from mcpclient import McpJujuClient

with McpJujuClient.stdio('../mcp-juju', args=['--server-type', 'stdio']) as mcp:
    print(len(mcp.list_tools()), 'tools')
    result = mcp.call_tool('status', {'model': 'mcp-juju-demo'})
    status = result.structured_content
    for name, unit in status['applications']['db']['units'].items():
        print(name, unit['workload-status']['current'], 'on machine', unit['machine'])
EOF
```

```
151 tools
db/0 active on machine 0
```

Notice two things. The tool is called `status`, exactly like the Juju command, and its flag `--model` became the property `model`. And the result arrived as `structured_content`, a parsed JSON object, because tools with a `--format` flag default to `json`.

## Step 5: Read a tool's documentation

Every tool also publishes the full Juju help text as a resource named `juju://<tool>-doc`. Still in `demo/`, read the one for `status`:

```bash
uv run python - <<'EOF'
from mcpclient import McpJujuClient

with McpJujuClient.stdio('../mcp-juju', args=['--server-type', 'stdio']) as mcp:
    doc = mcp.read_resource('juju://status-doc').contents[0].text
    print(doc[:400])
EOF
```

The output begins with `# status` and the Juju usage text:

```
# status

Report the status of the model, its machines, applications and units.

## Usage
juju status [options] [<selector> [...]]
...
```

This is the same text `juju help status` prints, so an assistant can look up any flag before using it.

## Step 6: See what read-only mode refuses

Start the server again with `--read-only`. Fewer tools are registered, and the ones that stay only accept queries. Try to change a configuration value:

```bash
uv run python - <<'EOF'
from mcpclient import McpJujuClient, result_text

args = ['--server-type', 'stdio', '--read-only']
with McpJujuClient.stdio('../mcp-juju', args=args) as ro:
    print(len(ro.list_tools()), 'tools')
    query = ro.call_tool('config', {'model': 'mcp-juju-demo', 'args': ['db']})
    print('query ok:', not query.is_error)
    write = ro.call_tool('config', {'model': 'mcp-juju-demo', 'args': ['db', 'profile=testing']})
    print('write is_error:', write.is_error)
    print(result_text(write))
EOF
```

```
60 tools
query ok: True
write is_error: True
rejected in read-only mode: setting "profile=testing" would modify configuration
```

The query went through, and the write was refused by the server before Juju was called. Confirm nothing changed:

```bash
juju config -m mcp-juju-demo db profile
```

```
production
```

Go back to the repository root before continuing:

```bash
cd ..
```

## Step 7: Run the full verification

The verify script repeats what you just did by hand, plus more checks, against three server configurations: stdio, `--read-only`, and Streamable HTTP with a bearer token. It prints what each tool returned and a verdict per check:

```bash
just demo-verify
```

The output ends with:

```
=== Summary ===
15 passed, 0 failed
```

Scroll up to see, among other things, the `models` list, the `config` settings of `db`, the exact rejection messages from read-only mode, and the `HTTP 401` returned when a request carries no token.

## Step 8: Browse the server in the MCP Inspector

The MCP Inspector is the reference tool for looking at any MCP server: it shows the tools, their schemas and resources, and lets you call them from a web page. Start it on the server in read-only mode:

```bash
just inspect --read-only
```

This runs `npx @modelcontextprotocol/inspector ./mcp-juju -- --server-type stdio --read-only`. Before launching, the recipe prints example calls for the demo model, each as one JSON object you can paste, and the same calls as `just call` lines for the terminal. The first run downloads the Inspector; it then prints the URL of the web UI (`http://localhost:6274` with a session token) and opens it in your browser.

In the browser:

1. Click **Connect**. The status turns to *Connected* and the server instructions appear.
2. Open the **Tools** tab. Sixty tools are listed; `deploy` and `destroy-model` are not among them because the server runs read-only.
3. Select `status`, switch the arguments to **Edit as JSON**, paste the first example and click **Execute Tool**:

   ```json
   {"model": "mcp-juju-demo"}
   ```

   The result panel shows the same JSON you printed in step 4, with `db/0` active.
4. Select `config` and execute this JSON. The result is marked as an error and reads `rejected in read-only mode: setting "profile=testing" would modify configuration`:

   ```json
   {"model": "mcp-juju-demo", "args": ["db", "profile=testing"]}
   ```

5. Open the **Resources** tab and select `status-doc`. The help text from step 5 appears in the panel.

Try a few more from the printed list, for example:

| Tool | Arguments | What you see |
|---|---|---|
| `show-unit` | `{"model": "mcp-juju-demo", "args": ["db/0"]}` | The unit's machine, address, open ports and relations |
| `config` | `{"model": "mcp-juju-demo", "args": ["db", "profile"]}` | Only the `profile` setting |
| `remove-unit` | `{"model": "mcp-juju-demo", "args": ["db/0"], "dry-run": true}` | Juju's report of what it would remove; the unit stays |
| `remove-unit` | `{"model": "mcp-juju-demo", "args": ["db/0"]}` | `rejected in read-only mode: --dry-run=true is required` |

The same calls work without the browser, through the Inspector's CLI mode. Stop the Inspector with `Ctrl+C` and run one:

```bash
just call show-unit '{"model": "mcp-juju-demo", "args": ["db/0"]}' --read-only
```

The tool result is printed as JSON, with the unit details in the first `content` block.

## Step 9: Clean up

Destroy the demo model. This too goes through the server, using its `destroy-model` tool:

```bash
just demo-clean
```

```
=== Model mcp-juju-demo ===
[   0.9s] mcp: destroy-model {'args': ['mcp-juju-demo'], 'no-prompt': True, 'destroy-storage': True, 'force': True, 'timeout': '600s'}
[  41.3s] model destroyed
```

Keep the `lxd` controller for the next tutorial, or remove it with `juju destroy-controller lxd --destroy-all-models --no-prompt`.

## What you've learned

- Every `mcp-juju` tool is one Juju command with the same name; positional arguments go in `args` and flags become properties.
- Tools with `--format` return parsed JSON as `structuredContent`, and each tool's help is available as the `juju://<tool>-doc` resource.
- `--read-only` is enforced by the server: queries work, writes come back as `isError` results and never reach Juju.
- The demo scripts (`just demo-deploy`, `just demo-verify`, `just demo-clean`) drive the server the same way an assistant does, with Jubilant checking the result.
- The MCP Inspector shows exactly what an assistant is offered: tool schemas, annotations, results and resources. `just inspect` opens it and `just call` runs one tool from the terminal.

Next steps:

- [Command reference](../reference/commands.md) lists all 151 tools with their annotations and read-only availability.
- [Read-only mode reference](../reference/read-only-mode.md) describes every policy and rejection message.
- [MCP interface reference](../reference/mcp-interface.md) and [CLI reference](../reference/cli.md) cover the transports, schemas and flags in full.
- [demo/README.md](../../demo/README.md) shows how to run the demo against another controller or charm.
