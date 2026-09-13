# Read-only mode reference

`--read-only` (`MCP_JUJU_READ_ONLY=true`) restricts the server to invocations that change neither Juju state (controller, model, client store) nor the host running the server. The restriction is enforced by the server before a Juju command is constructed; it does not rely on the client honouring tool annotations.

## Registration rules

| Command class | Registered | Call handling |
|---|---|---|
| `readOnlyHint: true` in [Command reference](commands.md) | yes | Accepted, subject to the host-write flag rule |
| Listed in the policy table below | yes | Accepted only in the described shape |
| Any other command | no | A direct `tools/call` is rejected with `rejected in read-only mode: <name> modifies state` |

`--tool-names` is applied first; names that are not allowed in read-only mode are dropped with a warning.

## Host-write flags

Rejected on every command, including read-only ones:

| Flag | Effect it would have |
|---|---|
| `output` | Writes the command output to a file |
| `filename` | Writes a bundle or backup archive to a file |
| `filepath` | Writes a downloaded charm to a file |
| `browser` | Launches the host's web browser |

Rejection message: `rejected in read-only mode: flag --<flag> writes to the host`.

## Policy table

Each policy is deny-by-default for flags: a set flag must be one of the generic flags, one of the policy's allowed flags, or one of its required flags.

Generic flags accepted by every policy: `model`, `controller`, `format`, `color`, `no-color`, `no-browser-login`.

| Tool | Positional arguments | Allowed flags | Required flags | Verified read path in Juju 3.6 |
|---|---|---|---|---|
| `config` | any number; none may contain `=` | — | — | `ConfigCommandBase.Init` records `GetOne`/`GetAll` only; `Run` calls `Application.Get` |
| `model-config` | any number; none may contain `=` | — | — | as above; `Run` calls `ModelGetWithMetadata` |
| `model-defaults` | any number; none may contain `=` | `cloud`, `region` | — | as above; `Run` calls `ModelDefaults` |
| `controller-config` | any number; none may contain `=` | — | — | as above; `Run` calls `ControllerConfig` |
| `application-storage` | any number; none may contain `=` | — | — | as above; `Run` calls `GetApplicationStorage` |
| `default-region` | at most 1 | — | — | `Run` returns after printing when no region and no `--reset` is given, before `UpdateCredential` |
| `default-credential` | at most 1 | — | — | `Run` returns after printing when no credential and no `--reset` is given, before `UpdateCredential` |
| `remove-application` | any number | — | `dry-run` = `true` | Client returns after `performDryRun`; facade `DestroyApplication` returns before `ApplyOperation` |
| `remove-unit` | any number | — | `dry-run` = `true` | Client returns after `performDryRun`; facade `DestroyUnit` returns before `ApplyOperation`. Not supported for Kubernetes units (Juju returns an error) |
| `remove-machine` | any number | — | `dry-run` = `true` | Client returns after `performDryRun`; facade `destroyMachine` returns before `Destroy` |

`deploy --dry-run` is not offered: the deployer ignores the flag for local charms.

## Rejection messages

All messages start with `rejected in read-only mode:` and are returned as `isError` tool results.

| Situation | Message suffix |
|---|---|
| Command has no policy | `<name> modifies state` |
| Host-write flag set | `flag --<flag> writes to the host` |
| Too many positional arguments | `at most <n> positional argument(s) are allowed for a query` |
| Positional argument contains `=` | `setting "<arg>" would modify configuration` |
| Required flag missing or different | `--<flag>=<value> is required` |
| Flag not allowed by the policy | `flag --<flag> is not allowed for a query` |

## Effect on tool definitions

In read-only mode:

| Field | Change |
|---|---|
| `annotations.readOnlyHint` | `true` for every registered tool |
| `annotations.destructiveHint` | `false` for every registered tool |
| `annotations.idempotentHint` | `true` for every registered tool |
| `description` | Tools with a policy gain a sentence starting with `Read-only mode:` describing the accepted shape |
| `instructions` | Gains a paragraph describing the mode |

The documentation resources and the `juju://config/...` template are unchanged.
