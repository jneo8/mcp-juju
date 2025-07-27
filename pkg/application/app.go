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
	// a.mcpServer.AddTool(
	// 	mcp.NewTool("listControllers", mcp.WithDescription("List all juju controllers")),
	// 	gethandleListControllerTool(a.client),
	// )
	// a.mcpServer.AddTool(
	// 	mcp.NewTool(
	// 		"listModels",
	// 		mcp.WithDescription("List all juju models"),
	// 		mcp.WithString("controller"),
	// 	),
	// 	gethandleListModelTool(a.client),
	// )
	// a.mcpServer.AddTool(
	// 	mcp.NewTool(
	// 		"getStatus",
	// 		mcp.WithDescription("Get juju status"),
	// 		mcp.WithString("controller"),
	// 		mcp.WithString("model"),
	// 		mcp.WithBoolean("includeStorage"),
	// 	),
	// 	gethandleGetStatusTool(a.client),
	// )
	// a.mcpServer.AddTool(
	// 	mcp.NewTool(
	// 		"getApplicationConfig",
	// 		mcp.WithDescription("Get juju application config"),
	// 		mcp.WithString("controller"),
	// 		mcp.WithString("model"),
	// 		mcp.WithString("application"),
	// 	),
	// 	gethandleGetApplicationConfigTool(a.client),
	// )
	// a.mcpServer.AddTool(
	// 	mcp.NewTool(
	// 		"setApplicationConfig",
	// 		mcp.WithDescription("Set juju application config"),
	// 		mcp.WithString("controller"),
	// 		mcp.WithString("model"),
	// 		mcp.WithString("application"),
	// 		mcp.WithString("key", mcp.Required()),
	// 		mcp.WithString("value", mcp.Required()),
	// 	),
	// 	gethandleSetApplicationConfigTool(a.client),
	// )
	// a.mcpServer.AddTool(
	// 	mcp.NewTool(
	// 		"addModel",
	// 		mcp.WithDescription("Adds a workload model"),
	// 		mcp.WithString("controller"),
	// 		mcp.WithString("model", mcp.Required()),
	// 		mcp.WithString("owner"),
	// 		mcp.WithString("config"),
	// 		mcp.WithString("credential"),
	// 		mcp.WithBoolean("NoSwitch"),
	// 		mcp.WithBoolean("CloudRegion"),
	// 	),
	// 	gethandleAddModelTool(a.client),
	// )
	return nil
}
