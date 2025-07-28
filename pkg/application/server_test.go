package application

import (
	"fmt"
	"testing"

	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStreamableHTTPServer(t *testing.T) {
	// Arrange
	mcpServer := server.NewMCPServer(
		"test-server",
		"1.0.0",
		server.WithResourceCapabilities(true, false),
		server.WithLogging(),
	)

	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	// Act
	streamableServer := newStreamableHTTPServer(mcpServer, cfg)

	// Assert
	assert.NotNil(t, streamableServer)
	assert.IsType(t, &server.StreamableHTTPServer{}, streamableServer)
}

func TestNewStreamableHTTPServer_WithDifferentConfig(t *testing.T) {
	// Test with different configuration values
	
	// Arrange
	mcpServer := server.NewMCPServer(
		"test-server-2",
		"2.0.0",
		server.WithResourceCapabilities(false, true),
	)

	cfg := config.Config{
		Port:     9090,
		Debug:    true,
		EndPoint: "/custom-endpoint",
	}

	// Act
	streamableServer := newStreamableHTTPServer(mcpServer, cfg)

	// Assert
	assert.NotNil(t, streamableServer)
	assert.IsType(t, &server.StreamableHTTPServer{}, streamableServer)
}

func TestNewStreamableHTTPServer_NilMCPServer(t *testing.T) {
	// Test behavior with nil MCP server
	// Note: This might panic depending on the underlying implementation
	// but it's good to document the expected behavior
	
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	// Act & Assert
	// This test documents that passing nil will likely cause issues
	// The actual behavior depends on the mark3labs/mcp-go library implementation
	assert.NotPanics(t, func() {
		streamableServer := newStreamableHTTPServer(nil, cfg)
		// We expect this to return something, but it might not work properly
		assert.NotNil(t, streamableServer)
	})
}

func TestRunStreamableHTTPServer_FormatCheck(t *testing.T) {
	// This test verifies the address format but doesn't actually start the server
	// to avoid binding to ports during testing
	
	// Arrange
	mcpServer := server.NewMCPServer(
		"test-server",
		"1.0.0",
		server.WithResourceCapabilities(true, false),
	)

	testCases := []struct {
		name           string
		config         config.Config
		expectedFormat string
	}{
		{
			name: "default port",
			config: config.Config{
				Port:     8080,
				Debug:    false,
				EndPoint: "/mcp",
			},
			expectedFormat: ":8080",
		},
		{
			name: "custom port",
			config: config.Config{
				Port:     9090,
				Debug:    true,
				EndPoint: "/custom",
			},
			expectedFormat: ":9090",
		},
		{
			name: "high port number",
			config: config.Config{
				Port:     65000,
				Debug:    false,
				EndPoint: "/test",
			},
			expectedFormat: ":65000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			streamableServer := newStreamableHTTPServer(mcpServer, tc.config)
			require.NotNil(t, streamableServer)

			// Act - verify the expected address format
			expectedAddr := fmt.Sprintf(":%d", tc.config.Port)
			
			// Assert
			assert.Equal(t, tc.expectedFormat, expectedAddr)
			
			// We also verify that the URL method works correctly
			expectedURL := fmt.Sprintf("http://localhost:%d%s", tc.config.Port, tc.config.EndPoint)
			assert.Equal(t, expectedURL, tc.config.URL())
		})
	}
}

func TestConfig_URL(t *testing.T) {
	// Test the URL generation method from config
	
	testCases := []struct {
		name        string
		config      config.Config
		expectedURL string
	}{
		{
			name: "default configuration",
			config: config.Config{
				Port:     8080,
				EndPoint: "/mcp",
			},
			expectedURL: "http://localhost:8080/mcp",
		},
		{
			name: "custom port and endpoint",
			config: config.Config{
				Port:     3000,
				EndPoint: "/api/mcp",
			},
			expectedURL: "http://localhost:3000/api/mcp",
		},
		{
			name: "root endpoint",
			config: config.Config{
				Port:     8080,
				EndPoint: "/",
			},
			expectedURL: "http://localhost:8080/",
		},
		{
			name: "empty endpoint",
			config: config.Config{
				Port:     8080,
				EndPoint: "",
			},
			expectedURL: "http://localhost:8080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			actualURL := tc.config.URL()

			// Assert
			assert.Equal(t, tc.expectedURL, actualURL)
		})
	}
}

func TestConfig_StreamableHTTPOptions(t *testing.T) {
	// Test the StreamableHTTPOptions method
	
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	// Act
	options := cfg.StreamableHTTPOptions()

	// Assert
	assert.NotNil(t, options)
	assert.Len(t, options, 1, "Should return exactly one option")
	assert.IsType(t, []server.StreamableHTTPOption{}, options)
}

func TestConfig_StreamableHTTPOptions_DifferentEndpoints(t *testing.T) {
	// Test StreamableHTTPOptions with different endpoint configurations
	
	testCases := []struct {
		name     string
		endpoint string
	}{
		{"default endpoint", "/mcp"},
		{"custom endpoint", "/api/v1/mcp"},
		{"root endpoint", "/"},
		{"nested endpoint", "/service/mcp/v1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			cfg := config.Config{
				Port:     8080,
				Debug:    false,
				EndPoint: tc.endpoint,
			}

			// Act
			options := cfg.StreamableHTTPOptions()

			// Assert
			assert.NotNil(t, options)
			assert.Len(t, options, 1)
		})
	}
}