package config

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	// Host is the interface the HTTP server binds to. Defaults to loopback;
	// binding elsewhere requires AuthToken or AllowNoAuth.
	Host       string
	Port       int
	Debug      bool
	EndPoint   string
	ServerType string   `mapstructure:"server-type"`
	ToolNames  []string `mapstructure:"tool-names"`
	// ReadOnly exposes only commands that do not modify state; commands that
	// both read and write (config, model-config, ...) reject writes.
	ReadOnly bool `mapstructure:"read-only"`

	// AuthToken, when set, is the bearer token HTTP clients must present.
	AuthToken string `mapstructure:"auth-token"`
	// AllowNoAuth permits binding a non-loopback host without AuthToken.
	AllowNoAuth bool `mapstructure:"allow-no-auth"`
	// CORSOrigins lists browser origins allowed to call the HTTP endpoint.
	// Empty means CORS headers are not emitted.
	CORSOrigins []string `mapstructure:"cors-origins"`
	// TLSCert and TLSKey enable HTTPS when both are set.
	TLSCert string `mapstructure:"tls-cert"`
	TLSKey  string `mapstructure:"tls-key"`
}

// ListenAddr is the host:port the HTTP server listens on.
func (c *Config) ListenAddr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// URL is the address clients use to reach the HTTP endpoint.
func (c *Config) URL() string {
	scheme := "http"
	if c.TLSEnabled() {
		scheme = "https"
	}
	host := c.Host
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("%s://%s%s", scheme, net.JoinHostPort(host, strconv.Itoa(c.Port)), c.EndPoint)
}

// TLSEnabled reports whether HTTPS is configured.
func (c *Config) TLSEnabled() bool {
	return c.TLSCert != "" && c.TLSKey != ""
}

// IsLoopbackHost reports whether Host only accepts connections from this
// machine. An empty host binds every interface and is therefore not loopback.
func (c *Config) IsLoopbackHost() bool {
	if c.Host == "localhost" {
		return true
	}
	ip := net.ParseIP(c.Host)
	return ip != nil && ip.IsLoopback()
}

func (c *Config) StreamableHTTPOptions() []server.StreamableHTTPOption {
	opts := []server.StreamableHTTPOption{
		server.WithEndpointPath(c.EndPoint),
	}
	if len(c.CORSOrigins) > 0 {
		opts = append(opts, server.WithStreamableHTTPCORS(
			server.WithCORSAllowedOrigins(c.CORSOrigins...),
		))
	}
	return opts
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
	if !c.IsHTTPServer() {
		return nil
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", c.Port)
	}
	if !c.IsLoopbackHost() && c.AuthToken == "" && !c.AllowNoAuth {
		return fmt.Errorf("refusing to serve HTTP on non-loopback host %q without --auth-token; "+
			"set --auth-token, or --allow-no-auth to expose every Juju command unauthenticated", c.Host)
	}
	if (c.TLSCert == "") != (c.TLSKey == "") {
		return errors.New("--tls-cert and --tls-key must be set together")
	}
	return nil
}
