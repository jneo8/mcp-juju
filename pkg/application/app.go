package application

import (
	"github.com/jneo8/mcp-juju/config"
	"github.com/jneo8/mcp-juju/pkg/jujuclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Application interface {
	RunSSE() error
}

type application struct {
	mcpServer *server.MCPServer
	config    config.Config
	client    jujuclient.Client
}

func NewApplication(cfg config.Config, client jujuclient.Client) (Application, error) {
	app := &application{
		mcpServer: server.NewMCPServer(
			config.MCPServerName,
			config.Version,
			server.WithResourceCapabilities(true, false),
			server.WithLogging(),
		),
		config: cfg,
		client: client,
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
	a.mcpServer.AddTool(
		mcp.NewTool("listControllers", mcp.WithDescription("List all juju controllers")),
		gethandleListControllerTool(a.client),
	)
	a.mcpServer.AddTool(
		mcp.NewTool(
			"listModels",
			mcp.WithDescription("List all juju models"),
			mcp.WithString("controller"),
		),
		gethandleListModelTool(a.client),
	)
	a.mcpServer.AddTool(
		mcp.NewTool(
			"getStatus",
			mcp.WithDescription("Get juju status"),
			mcp.WithString("controller"),
			mcp.WithString("model"),
			mcp.WithBoolean("includeStorage"),
		),
		gethandleGetStatusTool(a.client),
	)
	a.mcpServer.AddTool(
		mcp.NewTool(
			"getApplicationConfig",
			mcp.WithDescription("Get juju application config"),
			mcp.WithString("controller"),
			mcp.WithString("model"),
			mcp.WithString("application"),
		),
		gethandleGetApplicationConfigTool(a.client),
	)
	a.mcpServer.AddTool(
		mcp.NewTool(
			"setApplicationConfig",
			mcp.WithDescription("Get juju application config"),
			mcp.WithString("controller"),
			mcp.WithString("model"),
			mcp.WithString("application"),
			mcp.WithObject("settings"),
		),
		gethandleSetApplicationConfigTool(a.client),
	)
	return nil
}
