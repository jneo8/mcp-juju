package application

import (
	"fmt"

	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

func newStreamableHTTPServer(mcpServer *server.MCPServer, cfg config.Config) *server.StreamableHTTPServer {
	return server.NewStreamableHTTPServer(
		mcpServer, cfg.StreamableHTTPOptions()...,
	)
}

func runStreamableHTTPServer(server *server.StreamableHTTPServer, cfg config.Config) error {
	log.Debug().Msgf("Run Streamable HTTP Server at %s", cfg.URL())
	return server.Start(fmt.Sprintf(":%d", cfg.Port))
}

func runStdioServer(mcpServer *server.MCPServer) error {
	log.Debug().Msg("Run Stdio Server")
	return server.ServeStdio(mcpServer)
}
