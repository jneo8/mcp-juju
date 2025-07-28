package application

import (
	"github.com/jneo8/mcp-juju/config"
	"github.com/jneo8/mcp-juju/pkg/jujuadapter"
	"github.com/mark3labs/mcp-go/server"
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
			server.WithResourceCapabilities(true, false),
			server.WithLogging(),
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
	streamableHTTPServer := newStreamableHTTPServer(a.mcpServer, a.config)
	return runStreamableHTTPServer(streamableHTTPServer, a.config)
}

func (a *application) init() error {
	toolNames := a.adapter.ToolNames()
	for _, toolName := range toolNames {
		tool, handlerFunc, err := a.adapter.GetTool(toolName)
		if err != nil {
			return err
		}
		a.mcpServer.AddTool(*tool, handlerFunc)
	}
	return nil
}
