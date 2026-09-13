package cmd

import (
	"os"
	"testing"

	"github.com/jneo8/mcp-juju/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	t.Run("should have correct command configuration", func(t *testing.T) {
		assert.Equal(t, config.AppName, rootCmd.Use)
		assert.Equal(t, "MCP Juju", rootCmd.Short)
		assert.NotNil(t, rootCmd.RunE)
		assert.NotNil(t, rootCmd.PersistentPreRunE)
	})

	t.Run("should have required flags", func(t *testing.T) {

		portFlag := rootCmd.Flag("port")
		require.NotNil(t, portFlag)
		assert.Equal(t, "8080", portFlag.DefValue)

		debugFlag := rootCmd.Flag("debug")
		require.NotNil(t, debugFlag)
		assert.Equal(t, "false", debugFlag.DefValue)
	})
}

func TestPersistentPreRun(t *testing.T) {
	t.Run("should bind flags and unmarshal config", func(t *testing.T) {
		// Create a test command with flags
		testCmd := &cobra.Command{
			Use: "test",
		}
		testCmd.Flags().String("port", "8080", "Port to server on")
		testCmd.Flags().Bool("debug", false, "Enable debug mode")
		testCmd.Flags().String("server-type", "stdio", "Server type (http or stdio)")

		// Reset viper for test isolation
		viper.Reset()

		// Test the function
		err := persistentPreRun(testCmd, []string{})
		assert.NoError(t, err)

		// Verify that viper was configured
		assert.Equal(t, config.EnvPrefix, viper.GetEnvPrefix())
	})

	t.Run("should handle environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("MCP_JUJU_PORT", "9090")
		os.Setenv("MCP_JUJU_DEBUG", "true")
		defer func() {
			os.Unsetenv("MCP_JUJU_PORT")
			os.Unsetenv("MCP_JUJU_DEBUG")
		}()

		// Create a test command with flags
		testCmd := &cobra.Command{
			Use: "test",
		}
		testCmd.Flags().String("port", "8080", "Port to server on")
		testCmd.Flags().Bool("debug", false, "Enable debug mode")
		testCmd.Flags().String("server-type", "stdio", "Server type (http or stdio)")

		// Reset viper for test isolation
		viper.Reset()

		// Test the function
		err := persistentPreRun(testCmd, []string{})
		assert.NoError(t, err)

		// Verify that environment variables were read
		assert.Equal(t, "9090", viper.GetString("port"))
		assert.Equal(t, "true", viper.GetString("debug"))
	})

	t.Run("should map dashed flag names to underscored env vars", func(t *testing.T) {
		t.Setenv("MCP_JUJU_SERVER_TYPE", "http")
		t.Setenv("MCP_JUJU_AUTH_TOKEN", "from-env")

		testCmd := &cobra.Command{Use: "test"}
		testCmd.Flags().String("server-type", "stdio", "")
		testCmd.Flags().String("host", "127.0.0.1", "")
		testCmd.Flags().Int("port", 8080, "")
		testCmd.Flags().String("auth-token", "", "")
		viper.Reset()

		require.NoError(t, persistentPreRun(testCmd, []string{}))
		assert.Equal(t, "http", cfg.ServerType)
		assert.Equal(t, "from-env", cfg.AuthToken)
	})

	t.Run("should reject non-loopback http host without token", func(t *testing.T) {
		testCmd := &cobra.Command{Use: "test"}
		testCmd.Flags().String("server-type", "http", "")
		testCmd.Flags().String("host", "0.0.0.0", "")
		testCmd.Flags().Int("port", 8080, "")
		viper.Reset()

		err := persistentPreRun(testCmd, []string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "auth-token")
	})

	t.Run("should return error when server type is invalid", func(t *testing.T) {
		testCmd := &cobra.Command{
			Use: "test",
		}
		testCmd.Flags().String("server-type", "grpc", "Server type (http or stdio)")

		// Reset viper for test isolation
		viper.Reset()

		err := persistentPreRun(testCmd, []string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server type")
	})
}

func TestExecute(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		// We can't easily test the Execute function without mocking os.Exit
		// But we can at least verify the function exists and doesn't panic when called
		assert.NotPanics(t, func() {
			// Don't actually call Execute() as it would try to run the application
			// Just verify the function exists
			assert.NotNil(t, Execute)
		})
	})
}

func TestConfigUnmarshal(t *testing.T) {
	t.Run("should unmarshal config correctly", func(t *testing.T) {
		// Set up viper with test values
		viper.Reset()
		viper.Set("port", 9090)
		viper.Set("debug", true)

		var testConfig config.Config
		err := viper.Unmarshal(&testConfig)
		assert.NoError(t, err)

		assert.Equal(t, 9090, testConfig.Port)
		assert.Equal(t, true, testConfig.Debug)
	})
}
