package application

import (
	"context"
	"encoding/json"

	"github.com/jneo8/mcp-juju/pkg/jujuclient"
	"github.com/mark3labs/mcp-go/mcp"
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
