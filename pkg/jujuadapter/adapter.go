package jujuadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	jujucmd "github.com/juju/cmd/v3"
	"github.com/juju/gnuflag"
	"github.com/juju/juju/juju"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// formatFlagName is the Juju flag that selects the output formatter.
const formatFlagName = "format"

// jsonFormat is the formatter the adapter prefers so results can be returned
// as structured content.
const jsonFormat = "json"

type Adapter interface {
	ToolNames() []string
	GetTool(name string) (*mcp.Tool, mcpserver.ToolHandlerFunc, error)
	ToolDocResourceNames() []string
	GetResource(name string) (*mcp.Resource, mcpserver.ResourceHandlerFunc, error)
	ResourceTemplateNames() []string
	GetResourceTemplate(name string) (*mcp.ResourceTemplate, mcpserver.ResourceTemplateHandlerFunc, error)
}

func NewAdapter(toolNames []string) (Adapter, error) {
	a := &adapter{
		factory:   &commandFactory{},
		toolNames: toolNames,
	}
	a.init()
	return a, nil
}

type adapter struct {
	factory   CommandFactory
	toolNames []string
}

func (a *adapter) ToolNames() []string {
	// If specific tool names are configured, use those
	if len(a.toolNames) > 0 {
		return a.toolNames
	}

	// Otherwise, return all available command IDs
	ids := GetAllCommandIDs()
	names := make([]string, len(ids))
	for i, id := range ids {
		names[i] = string(id)
	}
	return names
}

func (a *adapter) ToolDocResourceNames() []string {
	// Create documentation resources for each tool (1-to-1 mapping)
	toolNames := a.ToolNames()
	resourceNames := make([]string, len(toolNames))
	for i, toolName := range toolNames {
		resourceNames[i] = toolName + "-doc"
	}
	return resourceNames
}

func (a *adapter) ResourceTemplateNames() []string {
	return a.factory.GetResourceTemplateNames()
}

func (a *adapter) init() {
	// Initialize Juju environment exactly like the CLI does - once at startup
	if err := juju.InitJujuXDGDataHome(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize Juju environment")
	}
}

func (a *adapter) buildEnhancedDescription(cmd Command) string {
	info := cmd.Info()
	if info == nil {
		return cmd.ToolDescription()
	}

	var desc strings.Builder

	// Start with purpose (short description)
	if info.Purpose != "" {
		desc.WriteString(info.Purpose)
	}

	// Add argument info if available
	if info.Args != "" {
		desc.WriteString(fmt.Sprintf("\n\nArguments: %s", info.Args))
	}

	desc.WriteString(fmt.Sprintf("\n\nFull help: read the resource juju://%s-doc.", cmd.Name()))

	result := desc.String()
	if result == "" {
		return cmd.ToolDescription()
	}
	return result
}

func (a *adapter) buildDocumentationContent(cmd Command) string {
	info := cmd.Info()
	if info == nil {
		return "No detailed documentation available."
	}

	var content strings.Builder

	// Start with purpose
	if info.Purpose != "" {
		content.WriteString(fmt.Sprintf("# %s\n\n%s\n\n", cmd.Name(), info.Purpose))
	}

	// Add arguments
	if info.Args != "" {
		content.WriteString(fmt.Sprintf("## Arguments\n\n%s\n\n", info.Args))
	}

	// Add detailed documentation
	if info.Doc != "" && info.Doc != info.Purpose {
		content.WriteString(fmt.Sprintf("## Details\n\n%s\n", strings.TrimSpace(info.Doc)))
	}

	result := content.String()
	if result == "" {
		return "No detailed documentation available."
	}
	return result
}

func (a *adapter) GetTool(name string) (*mcp.Tool, mcpserver.ToolHandlerFunc, error) {
	cmd, err := a.factory.GetCommandByName(name)
	if err != nil {
		return nil, nil, err
	}

	toolOptions, err := a.flagSetToToolOptions(cmd)
	if err != nil {
		return nil, nil, err
	}
	title := "juju " + cmd.Name()
	allOptions := []mcp.ToolOption{
		mcp.WithDescription(a.buildEnhancedDescription(cmd)),
		mcp.WithToolAnnotation(hintsFor(JujuCommandID(name)).toolAnnotation(title)),
	}
	allOptions = append(allOptions, toolOptions...)
	tool := mcp.NewTool(cmd.Name(), allOptions...)
	tool.Title = title
	handlerFunc := a.getHandlerFunc(cmd)
	return &tool, handlerFunc, nil
}

// flagKind is the JSON schema shape chosen for a Juju flag.
type flagKind int

const (
	flagString flagKind = iota
	flagBool
	flagInteger
	flagNumber
	flagDuration
	flagFile
	flagStringList // cmd.StringsValue: one comma-separated value
	flagAppendList // cmd.AppendStringsValue: flag may be repeated
	flagFormat     // cmd formatterValue: --format with a fixed set of choices
)

// classifyFlag determines the schema shape for a flag from its Value type.
// Juju's concrete value types are mostly unexported, so after checking the
// exported ones and the gnuflag boolFlag interface it falls back to the type
// name.
func classifyFlag(flag *gnuflag.Flag) flagKind {
	if bf, ok := flag.Value.(interface{ IsBoolFlag() bool }); ok && bf.IsBoolFlag() {
		return flagBool
	}
	switch flag.Value.(type) {
	case *jujucmd.StringsValue:
		return flagStringList
	case *jujucmd.AppendStringsValue:
		return flagAppendList
	case *jujucmd.FileVar:
		return flagFile
	}
	typeName := reflect.TypeOf(flag.Value).String()
	switch {
	case strings.HasSuffix(typeName, "formatterValue"):
		return flagFormat
	case strings.HasSuffix(typeName, "durationValue"):
		return flagDuration
	case strings.HasSuffix(typeName, "intValue"), strings.HasSuffix(typeName, "int64Value"),
		strings.HasSuffix(typeName, "uintValue"), strings.HasSuffix(typeName, "uint64Value"):
		return flagInteger
	case strings.HasSuffix(typeName, "float64Value"):
		return flagNumber
	default:
		return flagString
	}
}

var formatChoicesRE = regexp.MustCompile(`\(([a-z0-9|-]+)\)`)

// formatChoices extracts the allowed formatter names from the --format flag
// usage text, which Juju renders as "Specify output format (json|tabular|yaml)".
func formatChoices(usage string) []string {
	m := formatChoicesRE.FindStringSubmatch(usage)
	if m == nil {
		return nil
	}
	choices := strings.Split(m[1], "|")
	sort.Strings(choices)
	return choices
}

// defaultFormat returns the formatter the adapter should select when the
// caller does not choose one: json when the command supports it, otherwise
// the command's own default.
func defaultFormat(flag *gnuflag.Flag) string {
	for _, c := range formatChoices(flag.Usage) {
		if c == jsonFormat {
			return jsonFormat
		}
	}
	return flag.DefValue
}

// schemaFlag is a flag selected for the tool schema after alias folding.
type schemaFlag struct {
	flag  *gnuflag.Flag
	kind  flagKind
	usage string
}

// collectSchemaFlags returns the flags to expose, in flag-set order. Juju
// registers short aliases (for example -B for --base) as separate flags that
// share one Value; only the longest name of each group is exposed so the
// schema is not cluttered with duplicates.
func collectSchemaFlags(flagSet *gnuflag.FlagSet) []schemaFlag {
	type group struct {
		order int
		flags []*gnuflag.Flag
	}
	groups := map[any]*group{}
	var ordered []*group

	flagSet.VisitAll(func(flag *gnuflag.Flag) {
		var key any = flag // no sharing unless the Value is a pointer
		if v := reflect.ValueOf(flag.Value); v.Kind() == reflect.Ptr {
			key = v.Pointer()
		}
		g, ok := groups[key]
		if !ok {
			g = &group{order: len(ordered)}
			groups[key] = g
			ordered = append(ordered, g)
		}
		g.flags = append(g.flags, flag)
	})

	result := make([]schemaFlag, 0, len(ordered))
	for _, g := range ordered {
		chosen := g.flags[0]
		usage := chosen.Usage
		for _, f := range g.flags[1:] {
			if len(f.Name) > len(chosen.Name) {
				chosen = f
			}
			if usage == "" {
				usage = f.Usage
			}
		}
		if chosen.Usage != "" {
			usage = chosen.Usage
		}
		result = append(result, schemaFlag{flag: chosen, kind: classifyFlag(chosen), usage: usage})
	}
	return result
}

func (a *adapter) flagSetToToolOptions(cmd Command) ([]mcp.ToolOption, error) {
	flagSet := gnuflag.NewFlagSet(cmd.Name(), gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	toolOptions := []mcp.ToolOption{}

	// Add positional arguments support
	toolOptions = append(toolOptions, mcp.WithArray("args",
		mcp.WithStringItems(),
		mcp.Description("Positional arguments for the command, in CLI order"),
	))

	for _, sf := range collectSchemaFlags(flagSet) {
		flag := sf.flag
		switch sf.kind {
		case flagBool:
			opts := []mcp.PropertyOption{mcp.Description(sf.usage)}
			if flag.DefValue == "true" || flag.DefValue == "false" {
				opts = append(opts, mcp.DefaultBool(flag.DefValue == "true"))
			}
			toolOptions = append(toolOptions, mcp.WithBoolean(flag.Name, opts...))
		case flagInteger, flagNumber:
			opts := []mcp.PropertyOption{mcp.Description(sf.usage)}
			if def, err := strconv.ParseFloat(flag.DefValue, 64); err == nil {
				opts = append(opts, mcp.DefaultNumber(def))
			}
			toolOptions = append(toolOptions, mcp.WithNumber(flag.Name, opts...))
		case flagDuration:
			opts := []mcp.PropertyOption{mcp.Description(sf.usage + " (duration such as 30s, 5m or 1h)")}
			if flag.DefValue != "" {
				opts = append(opts, mcp.DefaultString(flag.DefValue))
			}
			toolOptions = append(toolOptions, mcp.WithString(flag.Name, opts...))
		case flagFile:
			toolOptions = append(toolOptions, mcp.WithString(flag.Name,
				mcp.Description(sf.usage+" (path to a local file)"),
			))
		case flagStringList, flagAppendList:
			toolOptions = append(toolOptions, mcp.WithArray(flag.Name,
				mcp.WithStringItems(),
				mcp.Description(sf.usage),
			))
		case flagFormat:
			opts := []mcp.PropertyOption{
				mcp.Description(sf.usage),
				mcp.DefaultString(defaultFormat(flag)),
			}
			if choices := formatChoices(flag.Usage); len(choices) > 0 {
				opts = append(opts, mcp.Enum(choices...))
			}
			toolOptions = append(toolOptions, mcp.WithString(flag.Name, opts...))
		default:
			opts := []mcp.PropertyOption{mcp.Description(sf.usage)}
			if flag.DefValue != "" {
				opts = append(opts, mcp.DefaultString(flag.DefValue))
			}
			toolOptions = append(toolOptions, mcp.WithString(flag.Name, opts...))
		}
	}
	return toolOptions, nil
}

func (a *adapter) getHandlerFunc(cmd Command) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return a.run(cmd.Name(), ctx, req)
	}
}

// flagValueToStrings converts an MCP argument into the string values passed
// to the flag's Set method. Arrays map to one Set call per element for
// repeatable flags and to a single comma-joined value for list flags.
func flagValueToStrings(kind flagKind, value interface{}) ([]string, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case string:
		if v == "" {
			return nil, nil
		}
		return []string{v}, nil
	case bool:
		return []string{strconv.FormatBool(v)}, nil
	case float64:
		return []string{strconv.FormatFloat(v, 'f', -1, 64)}, nil
	case int:
		return []string{strconv.Itoa(v)}, nil
	case []string:
		return joinListValues(kind, v)
	case []interface{}:
		items := make([]string, 0, len(v))
		for _, item := range v {
			if item == nil {
				continue
			}
			items = append(items, fmt.Sprintf("%v", item))
		}
		return joinListValues(kind, items)
	default:
		return []string{fmt.Sprintf("%v", v)}, nil
	}
}

func joinListValues(kind flagKind, items []string) ([]string, error) {
	if len(items) == 0 {
		return nil, nil
	}
	switch kind {
	case flagAppendList:
		return items, nil
	case flagStringList:
		return []string{strings.Join(items, ",")}, nil
	default:
		return nil, fmt.Errorf("expected a single value but got a list")
	}
}

// executionResult is the outcome of running a Juju command.
type executionResult struct {
	Stdout string
	Stderr string
	// Format is the output formatter that was in effect, or "" when the
	// command has no --format flag.
	Format string
}

// Output returns stdout and stderr combined, as a human-readable transcript.
func (r executionResult) Output() string {
	output := r.Stdout
	if r.Stderr != "" {
		if output != "" {
			output += "\n"
		}
		output += r.Stderr
	}
	return output
}

func (a *adapter) executeCommand(ctx context.Context, config CommandExecutionConfig) (executionResult, error) {
	var result executionResult

	// Get the command
	cmd, err := a.factory.GetCommandByName(config.CommandName)
	if err != nil {
		return result, fmt.Errorf("failed to get command '%s': %w", config.CommandName, err)
	}

	// Set up the command flags
	flagSet := gnuflag.NewFlagSet(config.CommandName, gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	// Set fixed flags first
	for flagName, flagValue := range config.FixedFlags {
		flag := flagSet.Lookup(flagName)
		if flag != nil {
			if err := flag.Value.Set(flagValue); err != nil {
				return result, fmt.Errorf("failed to set fixed flag '%s': %w", flagName, err)
			}
		}
	}

	// Set flag values from the configuration
	for key, value := range config.FlagValues {
		flag := flagSet.Lookup(key)
		if flag == nil {
			// Skip unknown flags
			continue
		}
		values, err := flagValueToStrings(classifyFlag(flag), value)
		if err != nil {
			return result, fmt.Errorf("invalid value for flag '%s': %w", key, err)
		}
		for _, v := range values {
			if err := flag.Value.Set(v); err != nil {
				return result, fmt.Errorf("failed to set flag '%s': %w", key, err)
			}
		}
	}

	// Prefer JSON output unless the caller chose a format, so the result can
	// also be returned as structured content.
	if flag := flagSet.Lookup(formatFlagName); flag != nil && classifyFlag(flag) == flagFormat {
		_, fixed := config.FixedFlags[formatFlagName]
		provided, ok := config.FlagValues[formatFlagName]
		chosen := fixed || (ok && provided != nil && provided != "")
		if !chosen {
			if err := flag.Value.Set(defaultFormat(flag)); err != nil {
				return result, fmt.Errorf("failed to set default format: %w", err)
			}
		}
		result.Format = flag.Value.String()
	}

	// Parse the flags (this validates the flag values)
	if err := flagSet.Parse(false, []string{}); err != nil {
		return result, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Initialize the command with positional arguments
	if err := cmd.Init(config.Arguments); err != nil {
		return result, fmt.Errorf("failed to initialize command '%s': %w", config.CommandName, err)
	}

	// Execute the command
	result.Stdout, result.Stderr, err = cmd.RunWithOutput(ctx)
	if err != nil {
		return result, fmt.Errorf("command '%s' failed: %w", config.CommandName, err)
	}
	return result, nil
}

// CommandExecutionConfig holds all the configuration needed to execute a command
type CommandExecutionConfig struct {
	CommandName string
	FixedFlags  map[string]string
	Arguments   []string
	FlagValues  map[string]interface{}
}

func (a *adapter) run(name string, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debug().Msgf("req: %#v", req)

	// Extract positional arguments from MCP request
	var positionalArgs []string
	var flagValues map[string]interface{}

	arguments, ok := req.Params.Arguments.(map[string]interface{})
	if ok {
		flagValues = make(map[string]interface{})

		// Extract positional args
		if argsInterface, exists := arguments["args"]; exists {
			switch argsValue := argsInterface.(type) {
			case []interface{}:
				for _, arg := range argsValue {
					if str, ok := arg.(string); ok && str != "" {
						positionalArgs = append(positionalArgs, str)
					}
				}
			case []string:
				for _, str := range argsValue {
					if str != "" {
						positionalArgs = append(positionalArgs, str)
					}
				}
			}
		}
		log.Debug().Msgf("Extracted positional args (filtered): %v", positionalArgs)

		// Extract flag values
		for key, value := range arguments {
			if key != "args" {
				flagValues[key] = value
			}
		}
	}

	config := CommandExecutionConfig{
		CommandName: name,
		Arguments:   positionalArgs,
		FlagValues:  flagValues,
	}

	result, err := a.executeCommand(ctx, config)
	if err != nil {
		// Command execution failures are tool errors, not protocol errors: the
		// MCP spec asks for isError results so the model can self-correct.
		return toolErrorResult(name, err, result), nil
	}

	return toolSuccessResult(result), nil
}

// toolErrorResult builds an isError result carrying the failure reason and
// whatever the command printed before failing.
func toolErrorResult(name string, err error, result executionResult) *mcp.CallToolResult {
	var text strings.Builder
	text.WriteString(fmt.Sprintf("juju %s failed: %v", name, err))
	if output := strings.TrimSpace(result.Output()); output != "" {
		text.WriteString("\n\n")
		text.WriteString(output)
	}
	return mcp.NewToolResultError(text.String())
}

// toolSuccessResult returns the command output as text and, when the output
// is a JSON object, also as structured content.
func toolSuccessResult(result executionResult) *mcp.CallToolResult {
	output := result.Output()
	if result.Format == jsonFormat {
		var structured map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &structured); err == nil && structured != nil {
			return mcp.NewToolResultStructured(structured, output)
		}
	}
	return mcp.NewToolResultText(output)
}

func (a *adapter) GetResource(name string) (*mcp.Resource, mcpserver.ResourceHandlerFunc, error) {
	// Check if this is a documentation resource (ends with -doc)
	if strings.HasSuffix(name, "-doc") {
		toolName := strings.TrimSuffix(name, "-doc")
		cmd, err := a.factory.GetCommandByName(toolName)
		if err != nil {
			return nil, nil, err
		}

		// Create documentation resource
		uri := fmt.Sprintf("juju://%s", name)
		resource := mcp.NewResource(
			uri,
			name,
			mcp.WithResourceDescription(fmt.Sprintf("Documentation for %s command", toolName)),
			mcp.WithMIMEType("text/markdown"),
		)

		handlerFunc := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			content := a.buildDocumentationContent(cmd)
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			}, nil
		}

		return &resource, handlerFunc, nil
	}

	return nil, nil, fmt.Errorf("resource '%s' not found", name)
}

func (a *adapter) GetResourceTemplate(name string) (*mcp.ResourceTemplate, mcpserver.ResourceTemplateHandlerFunc, error) {
	configs := a.factory.GetResourceTemplateConfigs()
	config, exists := configs[name]
	if !exists {
		return nil, nil, fmt.Errorf("resource template '%s' not found", name)
	}

	template := mcp.NewResourceTemplate(
		config.URITemplate,
		config.Name,
		mcp.WithTemplateDescription(config.Description),
		mcp.WithTemplateMIMEType("application/json"),
	)

	handlerFunc := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return a.handleResourceTemplate(ctx, req, config)
	}

	return &template, handlerFunc, nil
}

func (a *adapter) handleResourceTemplate(ctx context.Context, req mcp.ReadResourceRequest, config ResourceTemplateConfig) ([]mcp.ResourceContents, error) {
	uri := req.Params.URI

	// Parse URI parameters using the template configuration
	uriParams, err := parseURIParameters(uri, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI parameters: %w", err)
	}

	// Prepare positional arguments from URI parameters
	var args []string
	maxArgIndex := -1
	argMap := make(map[int]string)

	for uriVar, argIndexStr := range config.URIToArgs {
		if value, exists := uriParams[uriVar]; exists && value != "" {
			argIndex, err := strconv.Atoi(argIndexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid argument index '%s' for URI parameter '%s'", argIndexStr, uriVar)
			}
			argMap[argIndex] = value
			if argIndex > maxArgIndex {
				maxArgIndex = argIndex
			}
		}
	}

	// Build args array in correct order
	if maxArgIndex >= 0 {
		args = make([]string, maxArgIndex+1)
		for i := 0; i <= maxArgIndex; i++ {
			if value, exists := argMap[i]; exists {
				args[i] = value
			} else {
				args[i] = "" // Fill gaps with empty strings
			}
		}
		// Remove trailing empty strings
		for len(args) > 0 && args[len(args)-1] == "" {
			args = args[:len(args)-1]
		}
	}

	// Prepare flag values from URI parameters
	flagValues := make(map[string]interface{})
	for uriVar, flagName := range config.URIToFlags {
		if value, exists := uriParams[uriVar]; exists && value != "" {
			flagValues[flagName] = value
		}
	}

	execConfig := CommandExecutionConfig{
		CommandName: config.CommandName,
		FixedFlags:  config.FixedFlags,
		Arguments:   args,
		FlagValues:  flagValues,
	}

	result, err := a.executeCommand(ctx, execConfig)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     result.Output(),
		},
	}, nil
}
