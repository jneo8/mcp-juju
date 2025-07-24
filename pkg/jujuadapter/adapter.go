package jujuadapter

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/juju/gnuflag"
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

func (a *adapter) init() {}

func (a *adapter) GetTool(name string) (*mcp.Tool, mcpserver.ToolHandlerFunc, error) {
	cmd, err := a.factory.GetCommand(name)
	if err != nil {
		return nil, nil, err
	}

	toolOptions, err := a.flagSetToToolOptions(cmd)
	if err != nil {
		return nil, nil, err
	}
	allOptions := []mcp.ToolOption{mcp.WithDescription(cmd.ToolDescription())}
	allOptions = append(allOptions, toolOptions...)
	tool := mcp.NewTool(cmd.Name(), allOptions...)
	handlerFunc := a.getHandlerFunc(cmd)
	return &tool, handlerFunc, nil
}

func (a *adapter) flagSetToToolOptions(cmd Command) ([]mcp.ToolOption, error) {
	flagSet := gnuflag.NewFlagSet(cmd.Name(), gnuflag.ContinueOnError)
	cmd.SetFlags(flagSet)

	toolOptions := []mcp.ToolOption{}
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

	// Convert MCP request arguments to flag values
	if err := a.setFlagsFromRequest(flagSet, req); err != nil {
		return nil, err
	}

	// Parse the flags (this validates the flag values)
	if err := flagSet.Parse(false, []string{}); err != nil {
		return nil, err
	}

	// Initialize the command with parsed arguments
	if err := cmd.Init(flagSet.Args()); err != nil {
		return nil, err
	}

	stdout, stderr, err := cmd.RunWithOutput(ctx)
	if err != nil {
		return nil, err
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
