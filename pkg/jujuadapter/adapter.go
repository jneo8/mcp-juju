package jujuadapter

import (
	"context"
	"errors"
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
		subCmds: make(map[string]Command, len(commandList)),
	}
	a.init()
	return a, nil
}

type adapter struct {
	subCmds map[string]Command
}

func (a *adapter) ToolNames() []string {
	return commandList
}

func (a *adapter) init() {
	for _, name := range commandList {
		cmd := NewCommand(name)
		log.Debug().Msgf("Register cmd %s", cmd.Name())
		a.subCmds[cmd.Name()] = cmd
	}
}

func (a *adapter) GetTool(name string) (*mcp.Tool, mcpserver.ToolHandlerFunc, error) {
	cmd, ok := a.subCmds[name]
	if !ok {
		return nil, nil, errors.New("Tool not exists")
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
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "success",
				},
			},
		}, nil
	}
}
