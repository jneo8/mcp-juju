package application

import (
	"context"
	"errors"
	"testing"

	"github.com/jneo8/mcp-juju/config"
	mockjujuadapter "github.com/jneo8/mcp-juju/mocks/jujuadapter"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApplication_Success(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	testTools := []string{"status", "deploy", "bootstrap"}
	mockAdapter.EXPECT().ToolNames().Return(testTools)

	// Mock GetTool calls for each tool
	for _, toolName := range testTools {
		tool := mcp.NewTool(toolName, mcp.WithDescription("Test tool"))
		handlerFunc := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return &mcp.CallToolResult{}, nil
		}
		mockAdapter.EXPECT().GetTool(toolName).Return(&tool, handlerFunc, nil)
	}

	// Act
	app, err := NewApplication(cfg, mockAdapter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, app)
	assert.IsType(t, &application{}, app)
}

func TestNewApplication_GetToolError(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	testTools := []string{"status", "deploy"}
	mockAdapter.EXPECT().ToolNames().Return(testTools)

	// First tool succeeds
	tool1 := mcp.NewTool("status", mcp.WithDescription("Test tool"))
	handlerFunc1 := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{}, nil
	}
	mockAdapter.EXPECT().GetTool("status").Return(&tool1, handlerFunc1, nil)

	// Second tool fails
	mockAdapter.EXPECT().GetTool("deploy").Return(nil, nil, errors.New("failed to get tool"))

	// Act
	app, err := NewApplication(cfg, mockAdapter)

	// Assert
	require.Error(t, err)
	assert.Nil(t, app)
	assert.Contains(t, err.Error(), "failed to get tool")
}

func TestNewApplication_EmptyToolNames(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	mockAdapter.EXPECT().ToolNames().Return([]string{})

	// Act
	app, err := NewApplication(cfg, mockAdapter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, app)
}

func TestApplication_RunServer_InterfaceExists(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    false,
		EndPoint: "/mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	mockAdapter.EXPECT().ToolNames().Return([]string{})

	app, err := NewApplication(cfg, mockAdapter)
	require.NoError(t, err)

	// Assert that RunServer method exists
	assert.NotNil(t, app.RunServer, "RunServer method should exist")
	
	// Note: We don't actually call RunServer here because it would start a real HTTP server
	// Integration tests should cover the actual server startup
}

func TestApplication_InitWithManyTools(t *testing.T) {
	// Test with a larger number of tools to simulate real usage
	
	// Arrange
	cfg := config.Config{
		Port:     8080,
		Debug:    true,
		EndPoint: "/mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	
	// Create a list of many tools similar to what the real adapter would have
	testTools := []string{
		"version", "bootstrap", "add-relation", "status", "deploy",
		"add-unit", "remove-unit", "config", "expose", "unexpose",
		"add-model", "destroy-model", "switch", "login", "logout",
	}
	
	mockAdapter.EXPECT().ToolNames().Return(testTools)

	for _, toolName := range testTools {
		tool := mcp.NewTool(toolName, mcp.WithDescription("Juju "+toolName+" command"))
		handlerFunc := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: "Mock response for " + toolName,
					},
				},
			}, nil
		}
		mockAdapter.EXPECT().GetTool(toolName).Return(&tool, handlerFunc, nil)
	}

	// Act
	app, err := NewApplication(cfg, mockAdapter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, app)
}

func TestApplication_ConfigValues(t *testing.T) {
	// Test that the application properly stores configuration values
	
	// Arrange
	cfg := config.Config{
		Port:     9090,
		Debug:    true,
		EndPoint: "/custom-mcp",
	}

	mockAdapter := mockjujuadapter.NewMockAdapter(t)
	mockAdapter.EXPECT().ToolNames().Return([]string{"version"})
	
	tool := mcp.NewTool("version", mcp.WithDescription("Version tool"))
	handlerFunc := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{}, nil
	}
	mockAdapter.EXPECT().GetTool("version").Return(&tool, handlerFunc, nil)

	// Act
	app, err := NewApplication(cfg, mockAdapter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, app)
	
	// Verify the application was created with the correct config
	appImpl := app.(*application)
	assert.Equal(t, cfg.Port, appImpl.config.Port)
	assert.Equal(t, cfg.Debug, appImpl.config.Debug)
	assert.Equal(t, cfg.EndPoint, appImpl.config.EndPoint)
	assert.Equal(t, mockAdapter, appImpl.adapter)
	assert.NotNil(t, appImpl.mcpServer)
}