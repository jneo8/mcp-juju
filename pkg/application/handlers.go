package application

import (
	"context"
	"encoding/json"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jneo8/mcp-juju/pkg/jujuclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func gethandleListControllerTool(client jujuclient.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		controllers, err := client.GetControllers()
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(controllers)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(jsonBytes),
				},
			},
		}, nil
	}
}

func gethandleListModelTool(client jujuclient.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args ListModelToolArgs
		mapstructure.Decode(req.Params.Arguments, &args)
		models, err := client.GetModels(args.Controller)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(models)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(jsonBytes),
				},
			},
		}, nil
	}
}

func gethandleGetStatusTool(client jujuclient.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args GetStatusToolArgs
		mapstructure.Decode(req.Params.Arguments, &args)
		status, err := client.GetStatus(ctx, args.Controller, args.Model, args.IncludeStorage)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(status)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(jsonBytes),
				},
			},
		}, nil
	}
}

func gethandleGetApplicationConfigTool(client jujuclient.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args GetApplicationConfigToolArgs
		mapstructure.Decode(req.Params.Arguments, &args)
		status, err := client.GetApplicationConfig(ctx, args.Controller, args.Model, args.Application)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(status)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(jsonBytes),
				},
			},
		}, nil
	}
}

func gethandleSetApplicationConfigTool(client jujuclient.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args SetApplicationConfigToolArgs
		mapstructure.Decode(req.Params.Arguments, &args)
		log.Debug().Msgf("%#v %#v", args, req.Params.Arguments)
		err := client.SetApplicationConfig(ctx, args.Controller, args.Model, args.Application, map[string]string{args.Key: args.Value})
		if err != nil {
			return nil, err
		}
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
