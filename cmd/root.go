package cmd

import (
	"fmt"
	"os"

	"github.com/jneo8/mcp-juju/config"
	"github.com/jneo8/mcp-juju/pkg/application"
	"github.com/jneo8/mcp-juju/pkg/jujuadapter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

func init() {
	rootCmd.Flags().String("port", "8080", "Port to server on")
	rootCmd.Flags().String("endpoint", "/mcp", "Endpoint path for the server")
	rootCmd.Flags().String("server-type", "http", "Server type (http or stdio)")
	rootCmd.Flags().Bool("debug", false, "Enable debug mode")
}

var rootCmd = &cobra.Command{
	Use:               config.AppName,
	RunE:              run,
	Short:             "MCP Juju",
	PersistentPreRunE: persistentPreRun,
}

func run(cmd *cobra.Command, args []string) error {

	adapter, err := jujuadapter.NewAdapter()
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
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind flags: %w", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode config")
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
