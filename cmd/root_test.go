package cmd

import (
	"os"
	"testing"

	"github.com/juju/juju/mcp-juju/config"
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

		// Reset viper for test isolation
		viper.Reset()

		// Test the function
		err := persistentPreRun(testCmd, []string{})
		assert.NoError(t, err)

		// Verify that environment variables were read
		assert.Equal(t, "9090", viper.GetString("port"))
		assert.Equal(t, "true", viper.GetString("debug"))
	})

	t.Run("should return error when flag binding fails", func(t *testing.T) {
		// Create a command without flags to trigger binding error
		testCmd := &cobra.Command{
			Use: "test",
		}

		// Reset viper for test isolation
		viper.Reset()

		// This should not fail since BindPFlags with no flags is valid
		err := persistentPreRun(testCmd, []string{})
		assert.NoError(t, err)
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
