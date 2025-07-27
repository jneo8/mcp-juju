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
}

func NewAdapter() (Adapter, error) {
	a := &adapter{
		factory: &commandFactory{},
	}
	a.init()
	return a, nil
}

type adapter struct {
	factory CommandFactory
}

func (a *adapter) ToolNames() []string {
	return commandList
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

	// Add long documentation if available and different from purpose
	if info.Doc != "" && info.Doc != info.Purpose {
		desc.WriteString(fmt.Sprintf("\n\nDetails:\n%s", strings.TrimSpace(info.Doc)))
	}

	// Add examples if available
	if info.Examples != "" {
		desc.WriteString(fmt.Sprintf("\n\nExamples:\n%s", strings.TrimSpace(info.Examples)))
	}

	// Add see also if available
	if len(info.SeeAlso) > 0 {
		desc.WriteString(fmt.Sprintf("\n\nSee also: %s", strings.Join(info.SeeAlso, ", ")))
	}

	result := desc.String()
	if result == "" {
		return cmd.ToolDescription()
	}
	return result
}

func (a *adapter) GetTool(name string) (*mcp.Tool, mcpserver.ToolHandlerFunc, error) {
	cmd, err := a.factory.GetCommand(name)
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
			log.Debug().Msgf("flag: %#v", flag)

			// Convert flag to ToolOption based on its type
			// Use reflection to determine the type since concrete types may not be exported
			flagType := reflect.TypeOf(flag.Value).String()
			log.Debug().Msgf("flag type: %s", flagType)

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
				log.Debug().Msgf("Unknown or string flag type %s for flag %s, treating as string", flagType, flag.Name)
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

func (a *adapter) setFlagsFromRequest(flagSet *gnuflag.FlagSet, req mcp.CallToolRequest) error {
	// Iterate through the MCP request arguments and set flag values
	arguments, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid arguments type")
	}

	for key, value := range arguments {
		// Skip the special "args" parameter - it's handled separately
		if key == "args" {
			continue
		}

		flag := flagSet.Lookup(key)
		if flag == nil {
			// Skip unknown flags - they might be handled elsewhere
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
			return fmt.Errorf("failed to set flag %s: %w", key, err)
		}
	}
	return nil
}

func (a *adapter) run(name string, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debug().Msgf("req: %#v", req)
	cmd, err := a.factory.GetCommand(name)
	if err != nil {
		return nil, err
	}

	// Create flag set and populate it from the MCP request
	flagSet := gnuflag.NewFlagSet(cmd.Name(), gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	// Extract positional arguments from MCP request
	var positionalArgs []string
	arguments, ok := req.Params.Arguments.(map[string]interface{})
	if ok {
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
	}

	// Convert MCP request arguments to flag values
	if err := a.setFlagsFromRequest(flagSet, req); err != nil {
		return nil, err
	}

	// Parse the flags (this validates the flag values)
	if err := flagSet.Parse(false, []string{}); err != nil {
		return nil, err
	}

	// Initialize the command with positional arguments
	if err := cmd.Init(positionalArgs); err != nil {
		return nil, err
	}

	stdout, stderr, err := cmd.RunWithOutput(ctx)
	if err != nil {
		// Provide more context about the command that failed
		return nil, fmt.Errorf("command '%s' failed: %w\nStderr: %s", name, err, stderr)
	}

	// Combine stdout and stderr for the result
	output := stdout
	if stderr != "" {
		if output != "" {
			output += "\n"
		}
		output += stderr
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
