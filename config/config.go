package config

import (
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	Port       int
	Debug      bool
	EndPoint   string
	ServerType string   `mapstructure:"server-type"`
	ToolNames  []string `mapstructure:"tool-names"`
}

func (c *Config) URL() string {
	return fmt.Sprintf("http://localhost:%d%s", c.Port, c.EndPoint)
}

func (c *Config) StreamableHTTPOptions() []server.StreamableHTTPOption {
	return []server.StreamableHTTPOption{
		c.endpointPath(),
	}
}

func (c *Config) endpointPath() server.StreamableHTTPOption {
	return server.WithEndpointPath(c.EndPoint)
}

func (c *Config) IsHTTPServer() bool {
	return c.ServerType == ServerTypeHTTP
}

func (c *Config) IsStdioServer() bool {
	return c.ServerType == ServerTypeStdio
}

func (c *Config) Validate() error {
	if c.ServerType != ServerTypeHTTP && c.ServerType != ServerTypeStdio {
		return errors.New("invalid server type: must be 'http' or 'stdio'")
	}
	return nil
}
