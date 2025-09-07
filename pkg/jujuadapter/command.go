package jujuadapter

import (
	"bytes"
	"context"
	"time"

	"github.com/juju/cmd/v3"
	"github.com/juju/gnuflag"
)

// commandList is now generated from command definitions
// Use GetAllCommandIDs() to get the list of commands

type Command interface {
	// cmd.Command
	SetFlags(f *gnuflag.FlagSet)
	Init(args []string) error
	Name() string
	ToolDescription() string
	Info() *cmd.Info
	Run(context.Context) error
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


func (c *command) getContext(ctx context.Context) (*cmd.Context, error) {
	cmdCtx, err := cmd.DefaultContext()
	if err != nil {
		return nil, err
	}
	// Note: cmd/v3 Context might not have Context field
	return cmdCtx, nil
}

func (c *command) getContextWithOutput(ctx context.Context) (*cmd.Context, *bytes.Buffer, *bytes.Buffer, error) {
	cmdCtx, err := cmd.DefaultContext()
	if err != nil {
		return nil, nil, nil, err
	}
	// Note: cmd/v3 Context might not have Context field

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
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

func (c *command) Run(ctx context.Context) error {
	cmdCtx, err := c.getContext(ctx)
	if err != nil {
		return err
	}
	c.cmd.Run(cmdCtx)
	return nil
}

func (c *command) RunWithOutput(ctx context.Context) (string, string, error) {
	cmdCtx, stdout, stderr, err := c.getContextWithOutput(ctx)
	if err != nil {
		return "", "", err
	}

	err = c.cmd.Run(cmdCtx)
	return stdout.String(), stderr.String(), err
}
