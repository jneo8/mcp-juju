package jujuadapter

import (
	"fmt"
	"sort"
	"strings"
)

// readOnlyRejection prefixes every error produced by read-only mode.
// Functional tests and clients match on it.
const readOnlyRejection = "rejected in read-only mode"

// genericReadFlags only select the target or shape the output, so they are
// accepted on every command in read-only mode.
var genericReadFlags = []string{"model", "controller", "format", "color", "no-color", "no-browser-login"}

// hostWriteFlags make a command write to, or act on, the host that runs the
// server: --output and --filename/--filepath write files, --browser launches
// a program. They are rejected on every command in read-only mode, including
// commands that are otherwise read-only, because read-only covers the host
// as well as Juju state.
var hostWriteFlags = []string{"output", "filename", "filepath", "browser"}

// readOnlyPolicy describes the invocation shapes of a command that only
// read. Read-only mode is deny-by-default: a flag that is neither generic,
// allowed nor required is rejected, so a write flag added upstream cannot
// slip through.
//
// Every policy below was checked against the Juju 3.6 sources; see the
// comments on each entry for the code path that makes the shape read-only.
type readOnlyPolicy struct {
	// allowFlags may be set in addition to genericReadFlags.
	allowFlags []string
	// requireFlags must be set to exactly these values.
	requireFlags map[string]string
	// maxArgs caps the positional arguments; -1 means unlimited.
	maxArgs int
	// noKeyValue rejects positional arguments of the form key=value.
	noKeyValue bool
	// note is appended to the tool description in read-only mode.
	note string
}

// configPolicy covers commands built on cmd/juju/config.ConfigCommandBase.
// Its Init derives the actions from the arguments: key=value → SetArgs,
// --file → SetFile, --reset → Reset, bare keys → GetOne, nothing → GetAll.
// Run dispatches GetOne/GetAll to the get methods, which only call the read
// APIs (Application.Get, ModelGetWithMetadata, ModelDefaults,
// ControllerConfig, GetApplicationStorage).
func configPolicy(extraFlags ...string) readOnlyPolicy {
	return readOnlyPolicy{
		allowFlags: extraFlags,
		maxArgs:    -1,
		noKeyValue: true,
		note:       "only queries are accepted; key=value arguments, --reset and --file are rejected.",
	}
}

// dryRunPolicy covers remove-* commands whose --dry-run path returns before
// any state change: the client calls performDryRun and returns (or fails on
// controllers without dry-run support), and the controller facades
// (DestroyApplication, DestroyUnit, destroyMachine) return the computed
// info before ApplyOperation / Destroy.
//
// deploy also has --dry-run but is deliberately excluded: the deployer
// ignores the flag for local charms ("ignoring dry-run flag for local
// charms") and deploys anyway.
var dryRunPolicy = readOnlyPolicy{
	requireFlags: map[string]string{"dry-run": "true"},
	maxArgs:      -1,
	note:         "only --dry-run invocations are accepted; nothing is removed.",
}

// readOnlyPolicies lists every command that stays available in read-only
// mode although it is not annotated read-only. A unit test checks that
// every hintReadWrite command has an entry and that no entry is annotated
// read-only.
var readOnlyPolicies = map[JujuCommandID]readOnlyPolicy{
	CmdConfig:             configPolicy(),
	CmdModelConfig:        configPolicy(),
	CmdModelDefaults:      configPolicy("cloud", "region"), // filters for the read path
	CmdControllerConfig:   configPolicy(),
	CmdApplicationStorage: configPolicy(),
	// `default-region <cloud>` / `default-credential <cloud>`: Run returns
	// after printing when neither a second argument nor --reset is given,
	// before store.UpdateCredential.
	CmdDefaultRegion:     {maxArgs: 1, note: "only `default-region <cloud>` is accepted; a second argument or --reset is rejected."},
	CmdDefaultCredential: {maxArgs: 1, note: "only `default-credential <cloud>` is accepted; a second argument or --reset is rejected."},
	CmdRemoveApplication: dryRunPolicy,
	CmdRemoveUnit:        dryRunPolicy,
	CmdRemoveMachine:     dryRunPolicy,
}

// check returns an error describing the first thing that stops the
// invocation from being read-only, or nil when it only reads.
func (p readOnlyPolicy) check(args []string, flags map[string]interface{}) error {
	if p.maxArgs >= 0 && len(args) > p.maxArgs {
		return fmt.Errorf("%s: at most %d positional argument(s) are allowed for a query", readOnlyRejection, p.maxArgs)
	}
	if p.noKeyValue {
		for _, arg := range args {
			if strings.Contains(arg, "=") {
				return fmt.Errorf("%s: setting %q would modify configuration", readOnlyRejection, arg)
			}
		}
	}
	for name, want := range p.requireFlags {
		if got := flagValueString(flags[name]); got != want {
			return fmt.Errorf("%s: --%s=%s is required", readOnlyRejection, name, want)
		}
	}
	allowed := map[string]bool{}
	for _, name := range genericReadFlags {
		allowed[name] = true
	}
	for _, name := range p.allowFlags {
		allowed[name] = true
	}
	for name := range p.requireFlags {
		allowed[name] = true
	}
	var rejected []string
	for name, value := range flags {
		if flagIsSet(value) && !allowed[name] {
			rejected = append(rejected, name)
		}
	}
	if len(rejected) > 0 {
		sort.Strings(rejected)
		return fmt.Errorf("%s: flag --%s is not allowed for a query", readOnlyRejection, rejected[0])
	}
	return nil
}

// flagIsSet reports whether an MCP argument value would cause the flag to be
// applied, mirroring the empty-value skipping in flagValueToStrings.
func flagIsSet(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return false
	case string:
		return v != ""
	case bool:
		return v
	case []interface{}:
		return len(v) > 0
	case []string:
		return len(v) > 0
	default:
		return true
	}
}

// flagValueString renders a set flag value the way it would be passed to Juju;
// unset values render as "".
func flagValueString(value interface{}) string {
	if !flagIsSet(value) {
		return ""
	}
	switch v := value.(type) {
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// allowedInReadOnly reports whether a command may be registered in read-only
// mode: read-only commands always, others only with a policy.
func allowedInReadOnly(id JujuCommandID) bool {
	if hintsFor(id).readOnly {
		return true
	}
	_, ok := readOnlyPolicies[id]
	return ok
}

// readOnlyNote returns the description suffix for a command in read-only
// mode, or "" for commands that are read-only anyway.
func readOnlyNote(id JujuCommandID) string {
	if hintsFor(id).readOnly {
		return ""
	}
	if p, ok := readOnlyPolicies[id]; ok {
		return "Read-only mode: " + p.note
	}
	return ""
}

// checkReadOnly enforces read-only mode for one invocation.
func checkReadOnly(id JujuCommandID, args []string, flags map[string]interface{}) error {
	for _, name := range hostWriteFlags {
		if flagIsSet(flags[name]) {
			return fmt.Errorf("%s: flag --%s writes to the host", readOnlyRejection, name)
		}
	}
	if hintsFor(id).readOnly {
		return nil
	}
	policy, ok := readOnlyPolicies[id]
	if !ok {
		return fmt.Errorf("%s: %s modifies state", readOnlyRejection, id)
	}
	return policy.check(args, flags)
}
