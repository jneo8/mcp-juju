package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const initializeRequest = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`

func newTestHandler(t *testing.T, cfg config.Config) http.Handler {
	t.Helper()
	mcpServer := server.NewMCPServer("test-server", "1.0.0")
	return newHTTPHandler(newStreamableHTTPServer(mcpServer, cfg), cfg)
}

// postInitialize sends an MCP initialize request through the handler.
func postInitialize(handler http.Handler, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(initializeRequest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestHTTPHandler_NoTokenConfigured(t *testing.T) {
	cfg := config.Config{Host: "127.0.0.1", Port: 8080, EndPoint: "/mcp"}
	rec := postInitialize(newTestHandler(t, cfg), nil)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var resp struct {
		Result struct {
			ServerInfo struct{ Name string } `json:"serverInfo"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "test-server", resp.Result.ServerInfo.Name)
}

func TestHTTPHandler_BearerToken(t *testing.T) {
	cfg := config.Config{Host: "127.0.0.1", Port: 8080, EndPoint: "/mcp", AuthToken: "s3cret"}
	handler := newTestHandler(t, cfg)

	t.Run("missing header", func(t *testing.T) {
		rec := postInitialize(handler, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Header().Get("WWW-Authenticate"), "Bearer")
	})

	t.Run("wrong token", func(t *testing.T) {
		rec := postInitialize(handler, map[string]string{"Authorization": "Bearer nope"})
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("wrong scheme", func(t *testing.T) {
		rec := postInitialize(handler, map[string]string{"Authorization": "Basic s3cret"})
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("correct token", func(t *testing.T) {
		rec := postInitialize(handler, map[string]string{"Authorization": "Bearer s3cret"})
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	})
}

func TestHTTPHandler_UnknownPath(t *testing.T) {
	cfg := config.Config{Host: "127.0.0.1", Port: 8080, EndPoint: "/mcp"}
	req := httptest.NewRequest(http.MethodPost, "/other", strings.NewReader(initializeRequest))
	rec := httptest.NewRecorder()
	newTestHandler(t, cfg).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHTTPHandler_CORSPreflight(t *testing.T) {
	cfg := config.Config{Host: "127.0.0.1", Port: 8080, EndPoint: "/mcp", CORSOrigins: []string{"https://app.example.com"}}
	handler := newTestHandler(t, cfg)

	req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))

	req = httptest.NewRequest(http.MethodOptions, "/mcp", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}

// freePort asks the kernel for an unused loopback port.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestRunStreamableHTTPServer_ServesAndShutsDown(t *testing.T) {
	cfg := config.Config{Host: "127.0.0.1", Port: freePort(t), EndPoint: "/mcp", AuthToken: "tok"}
	mcpServer := server.NewMCPServer("test-server", "1.0.0")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runStreamableHTTPServer(ctx, mcpServer, cfg) }()

	// Wait for the listener, then check auth end to end.
	var resp *http.Response
	require.Eventually(t, func() bool {
		req, _ := http.NewRequest(http.MethodPost, cfg.URL(), strings.NewReader(initializeRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		req.Header.Set("Authorization", "Bearer tok")
		var err error
		resp, err = http.DefaultClient.Do(req)
		return err == nil
	}, 5*time.Second, 50*time.Millisecond)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	unauth, err := http.Post(cfg.URL(), "application/json", strings.NewReader(initializeRequest))
	require.NoError(t, err)
	unauth.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, unauth.StatusCode)

	cancel()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(shutdownTimeout + time.Second):
		t.Fatal("server did not shut down after context cancellation")
	}
	_, err = http.Get(fmt.Sprintf("http://%s/mcp", cfg.ListenAddr()))
	assert.Error(t, err, "listener should be closed")
}
