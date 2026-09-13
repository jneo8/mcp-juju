package application

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jneo8/mcp-juju/config"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

const (
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func newStreamableHTTPServer(mcpServer *server.MCPServer, cfg config.Config) *server.StreamableHTTPServer {
	return server.NewStreamableHTTPServer(
		mcpServer, cfg.StreamableHTTPOptions()...,
	)
}

// newHTTPHandler mounts the MCP endpoint behind the bearer-token check.
// mcp-go itself rejects loopback requests whose Host header is not localhost
// (DNS rebinding protection) and handles CORS when origins are configured.
func newHTTPHandler(streamable *server.StreamableHTTPServer, cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(cfg.EndPoint, withBearerAuth(cfg.AuthToken, streamable))
	return mux
}

// withBearerAuth requires "Authorization: Bearer <token>" on every request
// when token is non-empty. Tokens are compared in constant time.
func withBearerAuth(token string, next http.Handler) http.Handler {
	if token == "" {
		return next
	}
	expected := sha256.Sum256([]byte(token))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if strings.HasPrefix(header, prefix) {
			got := sha256.Sum256([]byte(strings.TrimPrefix(header, prefix)))
			if subtle.ConstantTimeCompare(expected[:], got[:]) == 1 {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Bearer realm="mcp-juju"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

// runStreamableHTTPServer serves until ctx is cancelled, then shuts down
// gracefully.
func runStreamableHTTPServer(ctx context.Context, mcpServer *server.MCPServer, cfg config.Config) error {
	streamable := newStreamableHTTPServer(mcpServer, cfg)
	httpServer := &http.Server{
		Addr:              cfg.ListenAddr(),
		Handler:           newHTTPHandler(streamable, cfg),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Info().
		Str("url", cfg.URL()).
		Bool("auth", cfg.AuthToken != "").
		Bool("tls", cfg.TLSEnabled()).
		Msg("Run Streamable HTTP Server")
	if !cfg.IsLoopbackHost() && cfg.AuthToken == "" {
		log.Warn().Str("host", cfg.Host).Msg("Serving every Juju command without authentication on a non-loopback host")
	}

	errCh := make(chan error, 1)
	go func() {
		var err error
		if cfg.TLSEnabled() {
			err = httpServer.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
		} else {
			err = httpServer.ListenAndServe()
		}
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info().Msg("Shutting down Streamable HTTP Server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if sErr := streamable.Shutdown(shutdownCtx); err == nil {
			err = sErr
		}
		<-errCh
		return err
	}
}

func runStdioServer(mcpServer *server.MCPServer) error {
	log.Debug().Msg("Run Stdio Server")
	return server.ServeStdio(mcpServer)
}
