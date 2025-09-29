package jujuadapter

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/juju/gnuflag"
	"github.com/juju/juju/juju"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

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
	allOptions := []mcp.ToolOption{mcp.WithDescription(a.buildEnhancedDescription(cmd))}
	allOptions = append(allOptions, toolOptions...)
	tool := mcp.NewTool(cmd.Name(), allOptions...)
	handlerFunc := a.getHandlerFunc(cmd)
	return &tool, handlerFunc, nil
}

func (a *adapter) flagSetToToolOptions(cmd Command) ([]mcp.ToolOption, error) {
	flagSet := gnuflag.NewFlagSet(cmd.Name(), gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	toolOptions := []mcp.ToolOption{}

	// Add positional arguments support
	toolOptions = append(toolOptions, mcp.WithArray("args",
		mcp.WithStringItems(),
		mcp.Description("Positional arguments for the command"),
	))

	flagSet.VisitAll(
		func(flag *gnuflag.Flag) {
			// Convert flag to ToolOption based on its type
			// Use reflection to determine the type since concrete types may not be exported
			flagType := reflect.TypeOf(flag.Value).String()
			switch {
			case strings.Contains(flagType, "boolValue") || strings.Contains(flagType, "Bool"):
				defaultBool := flag.DefValue == "true"
				toolOptions = append(toolOptions, mcp.WithBoolean(flag.Name,
					mcp.Description(flag.Usage),
					mcp.DefaultBool(defaultBool),
				))
			case strings.Contains(flagType, "intValue") || strings.Contains(flagType, "Int"):
				if defaultInt, err := strconv.ParseFloat(flag.DefValue, 64); err == nil {
					toolOptions = append(toolOptions, mcp.WithNumber(flag.Name,
						mcp.Description(flag.Usage),
						mcp.DefaultNumber(defaultInt),
					))
				} else {
					toolOptions = append(toolOptions, mcp.WithString(flag.Name,
						mcp.Description(flag.Usage),
						mcp.DefaultString(flag.DefValue),
					))
				}
			case strings.Contains(flagType, "float64Value") || strings.Contains(flagType, "Float64"):
				if defaultFloat, err := strconv.ParseFloat(flag.DefValue, 64); err == nil {
					toolOptions = append(toolOptions, mcp.WithNumber(flag.Name,
						mcp.Description(flag.Usage),
						mcp.DefaultNumber(defaultFloat),
					))
				} else {
					toolOptions = append(toolOptions, mcp.WithString(flag.Name,
						mcp.Description(flag.Usage),
						mcp.DefaultString(flag.DefValue),
					))
				}
			default:
				// For unknown types or string types, treat as string
				toolOptions = append(toolOptions, mcp.WithString(flag.Name,
					mcp.Description(flag.Usage),
					mcp.DefaultString(flag.DefValue),
				))
			}
		},
	)
	return toolOptions, nil
}

func (a *adapter) getHandlerFunc(cmd Command) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return a.run(cmd.Name(), ctx, req)
	}
}

func (a *adapter) executeCommand(ctx context.Context, config CommandExecutionConfig) (string, error) {
	// Get the command
	cmd, err := a.factory.GetCommandByName(config.CommandName)
	if err != nil {
		return "", fmt.Errorf("failed to get command '%s': %w", config.CommandName, err)
	}

	// Set up the command flags
	flagSet := gnuflag.NewFlagSet(config.CommandName, gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	// Set fixed flags first
	for flagName, flagValue := range config.FixedFlags {
		flag := flagSet.Lookup(flagName)
		if flag != nil {
			if err := flag.Value.Set(flagValue); err != nil {
				return "", fmt.Errorf("failed to set fixed flag '%s': %w", flagName, err)
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

		// Convert the value to string and set it
		var stringValue string
		switch v := value.(type) {
		case string:
			stringValue = v
		case bool:
			if v {
				stringValue = "true"
			} else {
				stringValue = "false"
			}
		case float64:
			stringValue = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			stringValue = strconv.Itoa(v)
		default:
			// For other types, convert to string
			stringValue = fmt.Sprintf("%v", v)
		}

		// Skip empty string values for flags
		if stringValue == "" {
			continue
		}

		if err := flag.Value.Set(stringValue); err != nil {
			return "", fmt.Errorf("failed to set flag '%s': %w", key, err)
		}
	}

	// Parse the flags (this validates the flag values)
	if err := flagSet.Parse(false, []string{}); err != nil {
		return "", fmt.Errorf("failed to parse flags: %w", err)
	}

	// Initialize the command with positional arguments
	if err := cmd.Init(config.Arguments); err != nil {
		return "", fmt.Errorf("failed to initialize command '%s': %w", config.CommandName, err)
	}

	// Execute the command
	stdout, stderr, err := cmd.RunWithOutput(ctx)
	if err != nil {
		return "", fmt.Errorf("command '%s' failed: %w\nStderr: %s", config.CommandName, err, stderr)
	}

	// Combine stdout and stderr for the result
	output := stdout
	if stderr != "" {
		if output != "" {
			output += "\n"
		}
		output += stderr
	}

	return output, nil
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

	output, err := a.executeCommand(ctx, config)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: output,
			},
		},
	}, nil
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

	output, err := a.executeCommand(ctx, execConfig)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     output,
		},
	}, nil
}

