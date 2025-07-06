package config

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	Port     int
	Debug    bool
	EndPoint string
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
