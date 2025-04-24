package application

import (
	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
)

type Application interface {
	RunSSE() error
}

type application struct {
	mcpServer *server.MCPServer
	config    config.Config
}

func NewApplication(cfg config.Config) (Application, error) {
	app := &application{
		mcpServer: server.NewMCPServer(
			config.MCPServerName,
			config.Version,
			server.WithResourceCapabilities(true, false),
			server.WithLogging(),
		),
		config: cfg,
	}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *application) RunSSE() error {
	sseServer := newSSEServer(a.mcpServer, a.config)
	return runSSE(sseServer, a.config)
}

func (a *application) init() error {
	return nil
}
