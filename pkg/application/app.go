package application

import (
	"github.com/jneo8/mcp-juju/config"
	"github.com/jneo8/mcp-juju/pkg/jujuadapter"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type Application interface {
	RunServer() error
}

type application struct {
	mcpServer *server.MCPServer
	config    config.Config
	adapter   jujuadapter.Adapter
}

func NewApplication(cfg config.Config, adapter jujuadapter.Adapter) (Application, error) {
	app := &application{
		mcpServer: server.NewMCPServer(
			config.MCPServerName,
			config.Version,
			// Tool and resource capabilities are registered implicitly when
			// the first tool or resource is added. Subscriptions and the
			// logging capability are intentionally not advertised: neither
			// is implemented, and both are deprecated or removed in the
			// 2026-07-28 MCP specification.
			server.WithInstructions(serverInstructions),
		),
		config:  cfg,
		adapter: adapter,
	}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *application) RunServer() error {
	if a.config.IsStdioServer() {
		return runStdioServer(a.mcpServer)
	}
	streamableHTTPServer := newStreamableHTTPServer(a.mcpServer, a.config)
	return runStreamableHTTPServer(streamableHTTPServer, a.config)
}

func (a *application) init() error {
	// Register tools
	toolNames := a.adapter.ToolNames()
	for _, toolName := range toolNames {
		log.Debug().Msgf("Register mcp tool %s", toolName)
		tool, handlerFunc, err := a.adapter.GetTool(toolName)
		if err != nil {
			return err
		}
		a.mcpServer.AddTool(*tool, handlerFunc)
	}

	// Register documentation resources
	docResourceNames := a.adapter.ToolDocResourceNames()
	for _, resourceName := range docResourceNames {
		log.Debug().Msgf("Register mcp resource %s", resourceName)
		resource, handlerFunc, err := a.adapter.GetResource(resourceName)
		if err != nil {
			return err
		}
		a.mcpServer.AddResource(*resource, handlerFunc)
	}

	// Register resource templates
	resourceTemplateNames := a.adapter.ResourceTemplateNames()
	for _, templateName := range resourceTemplateNames {
		log.Debug().Msgf("Register mcp resource template %s", templateName)
		template, handlerFunc, err := a.adapter.GetResourceTemplate(templateName)
		if err != nil {
			return err
		}
		a.mcpServer.AddResourceTemplate(*template, handlerFunc)
	}

	return nil
}
