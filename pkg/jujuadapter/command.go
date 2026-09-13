package jujuadapter

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/juju/cmd/v3"
	"github.com/juju/gnuflag"
)

// Command wraps a Juju cmd.Command so the adapter can introspect its flags
// and run it with captured output.
type Command interface {
	SetFlags(f *gnuflag.FlagSet)
	Init(args []string) error
	Name() string
	ToolDescription() string
	Info() *cmd.Info
	RunWithOutput(context.Context) (string, string, error)
}

type command struct {
	cmd  cmd.Command
	info *cmd.Info
	t    time.Time
}

func (c *command) Name() string {
	return c.info.Name
}

func (c *command) ToolDescription() string {
	return c.info.Purpose
}

func (c *command) Info() *cmd.Info {
	return c.info
}

// newCaptureContext builds a Juju command context whose stdout and stderr are
// captured and whose stdin is empty. Stdin must never be the process stdin:
// in stdio transport mode that is the MCP protocol stream, and a Juju command
// waiting for a confirmation prompt would swallow protocol messages.
func (c *command) newCaptureContext() (*cmd.Context, *bytes.Buffer, *bytes.Buffer, error) {
	cmdCtx, err := cmd.DefaultContext()
	if err != nil {
		return nil, nil, nil, err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmdCtx.Stdin = strings.NewReader("")
	cmdCtx.Stdout = stdout
	cmdCtx.Stderr = stderr

	return cmdCtx, stdout, stderr, nil
}

func (c *command) SetFlags(f *gnuflag.FlagSet) {
	c.cmd.SetFlags(f)
}

func (c *command) Init(args []string) error {
	return c.cmd.Init(args)
}

// RunWithOutput runs the command and returns its stdout, stderr and error.
// Juju's cmd.Context carries no Go context, so ctx cannot cancel a running
// command; it is accepted for interface symmetry with the MCP handler.
func (c *command) RunWithOutput(ctx context.Context) (string, string, error) {
	cmdCtx, stdout, stderr, err := c.newCaptureContext()
	if err != nil {
		return "", "", err
	}

	err = c.cmd.Run(cmdCtx)
	return stdout.String(), stderr.String(), err
}
