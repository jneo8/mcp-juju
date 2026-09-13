package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jneo8/mcp-juju/config"
	"github.com/jneo8/mcp-juju/pkg/application"
	"github.com/jneo8/mcp-juju/pkg/jujuadapter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

func init() {
	rootCmd.Flags().String("server-type", "stdio", "Server type (http or stdio)")
	rootCmd.Flags().Bool("debug", false, "Enable debug logging on stderr")
	rootCmd.Flags().StringSlice("tool-names", []string{}, "List of tool names to register (empty means all tools)")
	rootCmd.Flags().Bool("read-only", false, "Expose only tools that do not modify state; config-style tools reject writes")

	// HTTP server options.
	rootCmd.Flags().String("host", "127.0.0.1", "Interface to bind the HTTP server to (non-loopback hosts require --auth-token)")
	rootCmd.Flags().String("port", "8080", "Port for the HTTP server")
	rootCmd.Flags().String("endpoint", "/mcp", "Endpoint path for the HTTP server")
	rootCmd.Flags().String("auth-token", "", "Bearer token HTTP clients must send (also MCP_JUJU_AUTH_TOKEN)")
	rootCmd.Flags().Bool("allow-no-auth", false, "Allow serving HTTP on a non-loopback host without --auth-token")
	rootCmd.Flags().StringSlice("cors-origins", []string{}, "Browser origins allowed by CORS (empty disables CORS headers)")
	rootCmd.Flags().String("tls-cert", "", "TLS certificate file; enables HTTPS together with --tls-key")
	rootCmd.Flags().String("tls-key", "", "TLS private key file")
}

var rootCmd = &cobra.Command{
	Use:               config.AppName,
	RunE:              run,
	Short:             "MCP Juju",
	PersistentPreRunE: persistentPreRun,
}

func run(cmd *cobra.Command, args []string) error {

	adapter, err := jujuadapter.NewAdapter(cfg.ToolNames, cfg.ReadOnly)
	if err != nil {
		return err
	}
	app, err := application.NewApplication(cfg, adapter)
	if err != nil {
		return err
	}
	if err := app.RunServer(); err != nil {
		return err
	}
	return nil
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(config.EnvPrefix)
	// Map flag names such as server-type to MCP_JUJU_SERVER_TYPE.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind flags: %w", err)
	}
	// Decode into a fresh struct so stale values never survive a re-run.
	var decoded config.Config
	if err := viper.Unmarshal(&decoded); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}
	cfg = decoded
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	configureLogging(cfg)
	return nil
}

// configureLogging sets the global log level from the config. Logs always go
// to stderr, so they never interfere with the MCP stdio transport on stdout.
func configureLogging(cfg config.Config) {
	level := zerolog.InfoLevel
	if cfg.Debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(os.Stderr)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
