package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandReferenceIsCurrent fails when docs/reference/commands.md no
// longer matches the adapter tables; run `just docs-commands` to update it.
func TestCommandReferenceIsCurrent(t *testing.T) {
	t.Setenv("JUJU_DATA", t.TempDir())
	want, err := Render()
	require.NoError(t, err)

	path := filepath.Join("..", "..", defaultOutput)
	got, err := os.ReadFile(path)
	require.NoError(t, err, "missing %s; run `just docs-commands`", path)
	assert.Equal(t, want, string(got), "%s is out of date; run `just docs-commands`", path)
}
