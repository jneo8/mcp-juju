package jujuadapter

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	jujucmd "github.com/juju/cmd/v3"
	"github.com/juju/gnuflag"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

// newTestAdapter builds an adapter against an empty Juju client store so the
// tests never touch the developer's real controllers.
func newTestAdapter(t *testing.T) *adapter {
	t.Helper()
	t.Setenv("JUJU_DATA", t.TempDir())
	a, err := NewAdapter(nil, false)
	require.NoError(t, err)
	return a.(*adapter)
}

func callTool(t *testing.T, a *adapter, name string, args map[string]interface{}) *mcp.CallToolResult {
	t.Helper()
	_, handler, err := a.GetTool(name)
	require.NoError(t, err)
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args
	result, err := handler(context.Background(), req)
	require.NoError(t, err, "tool failures must be reported as isError results, not Go errors")
	require.NotNil(t, result)
	return result
}

func resultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	require.NotEmpty(t, result.Content)
	text, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok, "first content block should be text")
	return text.Text
}

func TestCommandIDsMatchJujuNames(t *testing.T) {
	f := &commandFactory{}
	for _, id := range GetAllCommandIDs() {
		c, err := f.GetCommand(id)
		require.NoError(t, err, id)
		assert.Equal(t, string(id), c.Name(), "command ID must equal the Juju command name")
	}
}

func TestAnnotationsCoverAllCommands(t *testing.T) {
	for _, id := range GetAllCommandIDs() {
		_, ok := commandHints[id]
		assert.True(t, ok, "command %s has no annotation classification", id)
	}
	for id := range commandHints {
		_, err := (&commandFactory{}).GetCommand(id)
		assert.NoError(t, err, "annotation refers to unknown command %s", id)
	}
}

func TestGetTool_Annotations(t *testing.T) {
	a := newTestAdapter(t)

	status, _, err := a.GetTool("status")
	require.NoError(t, err)
	assert.Equal(t, "juju status", status.Title)
	assert.Equal(t, "juju status", status.Annotations.Title)
	require.NotNil(t, status.Annotations.ReadOnlyHint)
	assert.True(t, *status.Annotations.ReadOnlyHint)
	assert.False(t, *status.Annotations.DestructiveHint)
	assert.True(t, *status.Annotations.IdempotentHint)
	assert.False(t, *status.Annotations.OpenWorldHint)

	destroy, _, err := a.GetTool("destroy-model")
	require.NoError(t, err)
	assert.False(t, *destroy.Annotations.ReadOnlyHint)
	assert.True(t, *destroy.Annotations.DestructiveHint)
}

func TestGetTool_Schema(t *testing.T) {
	a := newTestAdapter(t)

	tool, _, err := a.GetTool("status")
	require.NoError(t, err)
	props := tool.InputSchema.Properties

	args, ok := props["args"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "array", args["type"])

	format, ok := props["format"].(map[string]any)
	require.True(t, ok, "status should expose --format")
	assert.Equal(t, "string", format["type"])
	assert.Equal(t, "json", format["default"])
	assert.ElementsMatch(t, []string{"json", "line", "oneline", "short", "summary", "tabular", "yaml"}, format["enum"])

	assert.Contains(t, tool.Description, "juju://status-doc")

	deploy, _, err := a.GetTool("deploy")
	require.NoError(t, err)
	dprops := deploy.InputSchema.Properties
	assert.Contains(t, dprops, "base")
	assert.NotContains(t, dprops, "B", "short alias sharing a value with --base should be folded")
	overlay, ok := dprops["overlay"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "array", overlay["type"], "repeatable flags are arrays")
}

func TestClassifyFlag(t *testing.T) {
	fs := gnuflag.NewFlagSet("test", gnuflag.ContinueOnError)
	var (
		b       bool
		i       int
		s       string
		list    jujucmd.StringsValue
		appends jujucmd.AppendStringsValue
		file    jujucmd.FileVar
	)
	fs.BoolVar(&b, "flag-bool", false, "")
	fs.IntVar(&i, "flag-int", 0, "")
	fs.StringVar(&s, "flag-string", "", "")
	fs.Var(&list, "flag-list", "")
	fs.Var(&appends, "flag-append", "")
	fs.Var(&file, "flag-file", "")

	expect := map[string]flagKind{
		"flag-bool":   flagBool,
		"flag-int":    flagInteger,
		"flag-string": flagString,
		"flag-list":   flagStringList,
		"flag-append": flagAppendList,
		"flag-file":   flagFile,
	}
	for name, kind := range expect {
		assert.Equal(t, kind, classifyFlag(fs.Lookup(name)), name)
	}
}

func TestFlagValueToStrings(t *testing.T) {
	values, err := flagValueToStrings(flagAppendList, []interface{}{"a", "b"})
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, values)

	values, err = flagValueToStrings(flagStringList, []interface{}{"a", "b"})
	require.NoError(t, err)
	assert.Equal(t, []string{"a,b"}, values)

	_, err = flagValueToStrings(flagString, []interface{}{"a", "b"})
	assert.Error(t, err)

	values, err = flagValueToStrings(flagBool, true)
	require.NoError(t, err)
	assert.Equal(t, []string{"true"}, values)

	values, err = flagValueToStrings(flagInteger, float64(3))
	require.NoError(t, err)
	assert.Equal(t, []string{"3"}, values)

	values, err = flagValueToStrings(flagString, "")
	require.NoError(t, err)
	assert.Nil(t, values, "empty strings are treated as unset")
}

func TestFormatChoices(t *testing.T) {
	assert.Equal(t, []string{"json", "tabular", "yaml"}, formatChoices("Specify output format (json|tabular|yaml)"))
	assert.Nil(t, formatChoices("no choices here"))
}

func TestRun_Success(t *testing.T) {
	a := newTestAdapter(t)
	// Public cloud metadata is compiled into the client, so this needs no
	// controller and does not depend on the host's LXD or kubeconfig.
	result := callTool(t, a, "regions", map[string]interface{}{"args": []interface{}{"aws"}, "client": true})
	assert.False(t, result.IsError, resultText(t, result))
	assert.Contains(t, resultText(t, result), "us-east-1")
}

func TestRun_StructuredContent(t *testing.T) {
	a := newTestAdapter(t)
	// "controllers" only reads the (empty) local client store.
	result := callTool(t, a, "controllers", nil)
	assert.False(t, result.IsError)
	structured, ok := result.StructuredContent.(map[string]interface{})
	require.True(t, ok, "json output should be returned as structured content")
	assert.Contains(t, structured, "controllers")

	var fromText map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(resultText(t, result)), &fromText),
		"first text block must be exactly the JSON output")
	assert.Equal(t, structured, fromText)
}

func TestRun_FailureIsToolError(t *testing.T) {
	a := newTestAdapter(t)
	result := callTool(t, a, "show-controller", map[string]interface{}{
		"args": []interface{}{"no-such-controller"},
	})
	assert.True(t, result.IsError)
	text := resultText(t, result)
	assert.Contains(t, text, "juju show-controller failed")
	assert.Contains(t, text, "no-such-controller")
}

func TestRun_BadFlagValueIsToolError(t *testing.T) {
	a := newTestAdapter(t)
	result := callTool(t, a, "status", map[string]interface{}{
		"format": "not-a-format",
	})
	assert.True(t, result.IsError)
	assert.Contains(t, resultText(t, result), "unknown format")
}

func TestExecuteCommand_DefaultsToJSON(t *testing.T) {
	a := newTestAdapter(t)
	result, err := a.executeCommand(context.Background(), CommandExecutionConfig{CommandName: "controllers"})
	require.NoError(t, err)
	assert.Equal(t, "json", result.Format)

	result, err = a.executeCommand(context.Background(), CommandExecutionConfig{
		CommandName: "controllers",
		FlagValues:  map[string]interface{}{"format": "yaml"},
	})
	require.NoError(t, err)
	assert.Equal(t, "yaml", result.Format)
}

// --- read-only mode ---------------------------------------------------------

func newReadOnlyAdapter(t *testing.T, toolNames ...string) *adapter {
	t.Helper()
	t.Setenv("JUJU_DATA", t.TempDir())
	a, err := NewAdapter(toolNames, true)
	require.NoError(t, err)
	return a.(*adapter)
}

func TestReadOnlyPoliciesMatchHints(t *testing.T) {
	for id, hints := range commandHints {
		_, hasPolicy := readOnlyPolicies[id]
		if hints.readWrite {
			assert.True(t, hasPolicy, "%s: hintReadWrite commands need a read-only policy", id)
		}
		if hints.readOnly {
			assert.False(t, hasPolicy, "%s: read-only commands need no policy", id)
		}
	}
	for id := range readOnlyPolicies {
		_, ok := commandHints[id]
		assert.True(t, ok, "policy for unknown command %s", id)
	}
}

func TestReadOnly_ToolNames(t *testing.T) {
	a := newReadOnlyAdapter(t)
	names := map[string]bool{}
	for _, n := range a.ToolNames() {
		names[n] = true
	}
	assert.True(t, names["status"])
	assert.True(t, names["config"], "read-write tools with a policy stay available")
	assert.True(t, names["remove-unit"], "dry-run capable tools stay available")
	assert.False(t, names["deploy"], "deploy ignores --dry-run for local charms")
	assert.False(t, names["download"], "its purpose is writing a host file")
	assert.False(t, names["download-backup"], "its purpose is writing a host file")
	assert.False(t, names["destroy-model"])
	for _, n := range a.ToolNames() {
		assert.True(t, allowedInReadOnly(JujuCommandID(n)), n)
	}

	explicit := newReadOnlyAdapter(t, "status", "deploy", "config")
	assert.ElementsMatch(t, []string{"status", "config"}, explicit.ToolNames())
}

func TestReadOnly_Annotations(t *testing.T) {
	a := newReadOnlyAdapter(t)
	tool, _, err := a.GetTool("config")
	require.NoError(t, err)
	assert.True(t, *tool.Annotations.ReadOnlyHint)
	assert.False(t, *tool.Annotations.DestructiveHint)
	assert.Contains(t, tool.Description, "Read-only mode: only queries")

	remove, _, err := a.GetTool("remove-unit")
	require.NoError(t, err)
	assert.True(t, *remove.Annotations.ReadOnlyHint)
	assert.Contains(t, remove.Description, "only --dry-run")

	status, _, err := a.GetTool("status")
	require.NoError(t, err)
	assert.NotContains(t, status.Description, "Read-only mode")
}

func TestReadOnlyPolicy_Check(t *testing.T) {
	cfg := readOnlyPolicies[CmdConfig]
	assert.NoError(t, cfg.check([]string{"app"}, nil))
	assert.NoError(t, cfg.check([]string{"app", "profile"}, map[string]interface{}{"format": "json", "model": "m", "reset": ""}))
	assert.ErrorContains(t, cfg.check([]string{"app", "profile=testing"}, nil), "profile=testing")
	assert.ErrorContains(t, cfg.check([]string{"app"}, map[string]interface{}{"reset": []interface{}{"profile"}}), "--reset")
	assert.ErrorContains(t, cfg.check([]string{"app"}, map[string]interface{}{"file": "/tmp/x.yaml"}), "--file")
	assert.ErrorContains(t, cfg.check([]string{"app"}, map[string]interface{}{"some-new-flag": true}), "--some-new-flag",
		"unknown flags are rejected by default")

	defaults := readOnlyPolicies[CmdModelDefaults]
	assert.NoError(t, defaults.check(nil, map[string]interface{}{"cloud": "aws", "region": "us-east-1"}))

	region := readOnlyPolicies[CmdDefaultRegion]
	assert.NoError(t, region.check([]string{"aws"}, nil))
	assert.ErrorContains(t, region.check([]string{"aws", "us-east-1"}, nil), "at most 1")
	assert.ErrorContains(t, region.check([]string{"aws"}, map[string]interface{}{"reset": true}), "--reset")

	remove := readOnlyPolicies[CmdRemoveUnit]
	assert.NoError(t, remove.check([]string{"app/0"}, map[string]interface{}{"dry-run": true, "model": "m"}))
	assert.ErrorContains(t, remove.check([]string{"app/0"}, nil), "--dry-run=true is required")
	assert.ErrorContains(t, remove.check([]string{"app/0"}, map[string]interface{}{"dry-run": false}), "--dry-run=true is required")
	assert.ErrorContains(t, remove.check([]string{"app/0"}, map[string]interface{}{"dry-run": true, "force": true}), "--force")
	assert.ErrorContains(t, remove.check([]string{"app/0"}, map[string]interface{}{"dry-run": true, "no-prompt": true}), "--no-prompt")
}

func TestReadOnly_RejectsWritesAsToolErrors(t *testing.T) {
	a := newReadOnlyAdapter(t)

	write := callTool(t, a, "config", map[string]interface{}{"args": []interface{}{"app", "k=v"}})
	assert.True(t, write.IsError)
	assert.Contains(t, resultText(t, write), readOnlyRejection)

	// A command that is not registered in read-only mode is still refused if
	// called directly.
	deploy := callTool(t, a, "deploy", map[string]interface{}{"args": []interface{}{"postgresql"}})
	assert.True(t, deploy.IsError)
	assert.Contains(t, resultText(t, deploy), readOnlyRejection)

	// Queries pass the guard (this one then fails for lack of a controller,
	// which is not a read-only rejection).
	query := callTool(t, a, "config", map[string]interface{}{"args": []interface{}{"app"}})
	assert.NotContains(t, resultText(t, query), readOnlyRejection)

	read := callTool(t, a, "regions", map[string]interface{}{"args": []interface{}{"aws"}, "client": true})
	assert.False(t, read.IsError, resultText(t, read))

	// Host writes are rejected even on read-only commands.
	out := callTool(t, a, "regions", map[string]interface{}{"args": []interface{}{"aws"}, "output": "/tmp/regions.json"})
	assert.True(t, out.IsError)
	assert.Contains(t, resultText(t, out), "--output writes to the host")
	assert.NoFileExists(t, "/tmp/regions.json")
}

func TestCheckReadOnly_HostWriteFlags(t *testing.T) {
	assert.NoError(t, checkReadOnly(CmdStatus, nil, map[string]interface{}{"format": "json"}))
	assert.ErrorContains(t, checkReadOnly(CmdStatus, nil, map[string]interface{}{"output": "x"}), "--output writes to the host")
	assert.ErrorContains(t, checkReadOnly(CmdExportBundle, nil, map[string]interface{}{"filename": "b.yaml"}), "--filename writes to the host")
	assert.ErrorContains(t, checkReadOnly(CmdDashboard, nil, map[string]interface{}{"browser": true}), "--browser writes to the host")
	assert.ErrorContains(t, checkReadOnly(CmdConfig, []string{"app"}, map[string]interface{}{"output": "x"}), "--output")
}
