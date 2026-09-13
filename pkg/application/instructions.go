package application

// serverInstructions is sent to clients during initialization and tells the
// model how the Juju tools are shaped.
const serverInstructions = `This server exposes the Juju 3.6 CLI. Each tool is one juju command and shares its name (for example "status", "deploy", "integrate").

How to call a tool:
- Put positional arguments in the "args" array, in the same order as on the command line.
- Every other property is a command-line flag; use the long flag name without dashes (for example "model", "no-prompt").
- Multi-value flags accept arrays. Flags never read from stdin, so pass a file path or an inline value instead.
- Commands that support --format default to "json"; the result then also carries structuredContent.
- Commands that normally ask for confirmation (destroy-model, remove-*, kill-controller) cannot prompt here. Pass "no-prompt": true or the equivalent "yes"/"force" flag the command documents.
- A failed command returns isError with the Juju error text; read it and adjust the call.

Documentation:
- Read the resource "juju://<tool>-doc" for the full help text of a command.
- The template "juju://config/{application}{/config_name*}" returns application configuration as JSON.

Tool annotations mark which commands are read-only and which are destructive. Confirm with the user before running destructive commands against production models.`
