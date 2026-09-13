package jujuadapter

import "strings"

// CommandDescription is the documentation view of one command: what
// tools/list advertises plus how read-only mode treats it. It feeds the
// generated command reference (tools/gendocs).
type CommandDescription struct {
	Name        string
	Purpose     string
	ReadOnly    bool
	Destructive bool
	Idempotent  bool
	OpenWorld   bool
	// ReadOnlyMode is "yes" (registered as is), "no" (not registered) or
	// "policy: <accepted shape>".
	ReadOnlyMode string
}

// DescribeCommands returns every command in registration order.
func DescribeCommands() ([]CommandDescription, error) {
	factory := &commandFactory{}
	ids := GetAllCommandIDs()
	out := make([]CommandDescription, 0, len(ids))
	for _, id := range ids {
		c, err := factory.GetCommand(id)
		if err != nil {
			return nil, err
		}
		h := hintsFor(id)
		mode := "no"
		if h.readOnly {
			mode = "yes"
		} else if p, ok := readOnlyPolicies[id]; ok {
			mode = "policy: " + strings.TrimSuffix(p.note, ".")
		}
		purpose := strings.TrimSpace(strings.SplitN(c.Info().Purpose, "\n", 2)[0])
		out = append(out, CommandDescription{
			Name:         string(id),
			Purpose:      purpose,
			ReadOnly:     h.readOnly,
			Destructive:  h.destructive,
			Idempotent:   h.idempotent,
			OpenWorld:    h.openWorld,
			ReadOnlyMode: mode,
		})
	}
	return out, nil
}

// UnexposedCommand is a Juju CLI command that is deliberately not a tool.
type UnexposedCommand struct {
	Name    string
	Purpose string
}

// UnexposedCommands lists the Juju 3.6 CLI commands that are not exposed
// because their constructors are unexported in
// github.com/juju/juju/cmd/juju/commands. Keep in sync with the TODO in
// command_defs.go.
var UnexposedCommands = []UnexposedCommand{
	{"bootstrap", "Initialise a cloud environment"},
	{"debug-log", "Display the log for a model"},
	{"enable-ha", "Ensure that sufficient controllers exist to provide redundancy"},
	{"help-action-commands", "Show help on a Juju charm action command"},
	{"help-hook-commands", "Show help on a Juju charm hook command"},
	{"migrate", "Migrate a model to another controller"},
	{"switch", "Select or identify the current controller and model"},
	{"sync-agent-binary", "Copy agent binaries from the official agent store into a local controller"},
	{"upgrade-controller", "Upgrade the Juju agent for a controller"},
	{"upgrade-model", "Upgrade the Juju agent for a model"},
	{"version", "Print the Juju CLI client version"},
}
