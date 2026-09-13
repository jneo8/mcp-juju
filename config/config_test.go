package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpConfig() Config {
	return Config{Host: "127.0.0.1", Port: 8080, EndPoint: "/mcp", ServerType: ServerTypeHTTP}
}

func TestValidate_ServerType(t *testing.T) {
	cfg := Config{ServerType: "grpc"}
	assert.ErrorContains(t, cfg.Validate(), "invalid server type")

	cfg = Config{ServerType: ServerTypeStdio}
	assert.NoError(t, cfg.Validate(), "stdio needs no HTTP settings")
}

func TestValidate_NonLoopbackRequiresAuth(t *testing.T) {
	for _, host := range []string{"0.0.0.0", "", "::", "10.0.0.5", "example.com"} {
		cfg := httpConfig()
		cfg.Host = host
		assert.ErrorContains(t, cfg.Validate(), "auth-token", "host %q", host)

		cfg.AuthToken = "secret"
		assert.NoError(t, cfg.Validate(), "host %q with token", host)

		cfg.AuthToken = ""
		cfg.AllowNoAuth = true
		assert.NoError(t, cfg.Validate(), "host %q with allow-no-auth", host)
	}
}

func TestValidate_LoopbackNeedsNoAuth(t *testing.T) {
	for _, host := range []string{"127.0.0.1", "::1", "localhost", "127.0.0.2"} {
		cfg := httpConfig()
		cfg.Host = host
		assert.NoError(t, cfg.Validate(), "host %q", host)
	}
}

func TestValidate_Port(t *testing.T) {
	for _, port := range []int{0, -1, 65536} {
		cfg := httpConfig()
		cfg.Port = port
		assert.ErrorContains(t, cfg.Validate(), "invalid port")
	}
}

func TestValidate_TLSPair(t *testing.T) {
	cfg := httpConfig()
	cfg.TLSCert = "cert.pem"
	assert.ErrorContains(t, cfg.Validate(), "tls-cert and --tls-key")

	cfg.TLSKey = "key.pem"
	require.NoError(t, cfg.Validate())
	assert.True(t, cfg.TLSEnabled())
}

func TestAddresses(t *testing.T) {
	cfg := httpConfig()
	assert.Equal(t, "127.0.0.1:8080", cfg.ListenAddr())
	assert.Equal(t, "http://127.0.0.1:8080/mcp", cfg.URL())

	cfg.Host = "::1"
	assert.Equal(t, "[::1]:8080", cfg.ListenAddr())
	assert.Equal(t, "http://[::1]:8080/mcp", cfg.URL())

	cfg.Host = ""
	assert.Equal(t, ":8080", cfg.ListenAddr())
	assert.Equal(t, "http://localhost:8080/mcp", cfg.URL())

	cfg.TLSCert, cfg.TLSKey = "c", "k"
	assert.Equal(t, "https://localhost:8080/mcp", cfg.URL())
}

func TestStreamableHTTPOptions(t *testing.T) {
	cfg := httpConfig()
	assert.Len(t, cfg.StreamableHTTPOptions(), 1, "endpoint only")

	cfg.CORSOrigins = []string{"https://app.example.com"}
	assert.Len(t, cfg.StreamableHTTPOptions(), 2, "endpoint and CORS")
}
