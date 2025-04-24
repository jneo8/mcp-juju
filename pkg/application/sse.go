package application

import (
	"fmt"

	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
)

func newSSEServer(mcpServer *server.MCPServer, cfg config.Config) *server.SSEServer {
	return server.NewSSEServer(mcpServer, server.WithBaseURL(cfg.URL()))
}

func runSSE(sseServer *server.SSEServer, cfg config.Config) error {
	return sseServer.Start(fmt.Sprintf(":%d", cfg.Port))
}
