package jujuadapter

import (
	"bytes"
	"context"
	"time"

	"github.com/juju/gnuflag"
	"github.com/juju/juju/internal/cmd"
)

var commandList = []string{
	// From registerCommands in exact order
	"version",
	// Creation commands.
	"bootstrap",
	"add-relation",

	// Cross model relations commands.
	"offer",
	"remove-offer",
	"show-offered-endpoint",
	"list-endpoints",
	"find-endpoints",
	"consume",
	"suspend-relation",
	"resume-relation",

	// Firewall rule commands.
	"set-firewall-rule",
	"list-firewall-rules",

	// Destruction commands.
	"remove-relation",
	"remove-application",
	"remove-unit",
	"remove-saas",

	// Reporting commands.
	"status",
	"switch",
	"status-history",

	// Error resolution and debugging commands.
	"exec",
	"scp",
	"ssh",
	"resolved",
	"debug-log",
	"debug-hooks",
	"debug-code",

	// Configuration commands.
	"get-constraints",
	"set-constraints",
	"sync-agent-binary",
	"upgrade-model",
	"upgrade-controller",
	"refresh",
	"bind",

	// Charm tool commands.
	"help-hooks",
	"help-actions",

	// Manage backups.
	"create-backup",
	"download-backup",

	// Manage authorized ssh keys.
	"add-ssh-key",
	"remove-ssh-key",
	"import-ssh-key",
	"ssh-keys",

	// Manage users and access
	"add-user",
	"change-password",
	"show-user",
	"users",
	"enable-user",
	"disable-user",
	"login",
	"logout",
	"remove-user",
	"whoami",

	// Manage machines
	"add-machine",
	"remove-machine",
	"machines",
	"show-machine",

	// Manage model
	"model-config",
	"model-defaults",
	"retry-provisioning",
	"destroy-model",
	"grant",
	"revoke",
	"show-model",
	"model-credential",

	"migrate",
	"export-bundle",

	// Manage and control actions
	"actions",
	"show-action",
	"cancel-action",
	"run",
	"operations",
	"show-operation",
	"show-task",

	// Manage and control applications
	"add-unit",
	"config",
	"deploy",
	"expose",
	"unexpose",
	"diff-bundle",
	"show-application",
	"show-unit",

	// Operation protection commands
	"disable-command",
	"disabled-commands",
	"enable-command",

	// Manage storage
	"add-storage",
	"storage",
	"create-storage-pool",
	"storage-pools",
	"remove-storage-pool",
	"update-storage-pool",
	"show-storage",
	"remove-storage",
	"detach-storage",
	"attach-storage",
	"import-filesystem",

	// Manage spaces
	"add-space",
	"spaces",
	"move-to-space",
	"reload-spaces",
	"show-space",
	"remove-space",
	"rename-space",

	// Manage subnets
	"subnets",

	// Manage controllers
	"add-model",
	"destroy-controller",
	"models",
	"kill-controller",
	"controllers",
	"register",
	"unregister",
	"enable-destroy-controller",
	"show-controller",
	"controller-config",

	// Manage clouds and credentials
	"update-cloud",
	"update-public-clouds",
	"clouds",
	"regions",
	"show-cloud",
	"add-cloud",
	"remove-cloud",
	"credentials",
	"detect-credentials",
	"set-default-region",
	"set-default-credential",
	"add-credential",
	"remove-credential",
	"update-credential",
	"show-credential",
	"grant-cloud",
	"revoke-cloud",

	// CAAS commands
	"add-k8s",
	"update-k8s",
	"remove-k8s",
	"scale-application",

	// Manage Application Credential Access
	"trust",

	// Juju Dashboard commands.
	"dashboard",

	// Resource commands.
	"attach-resource",
	"resources",
	"charm-resources",

	// CharmHub related commands
	"info",
	"find",
	"download",

	// Secrets.
	"secrets",
	"show-secret",
	"add-secret",
	"update-secret",
	"remove-secret",
	"grant-secret",
	"revoke-secret",

	// Secret backends.
	"secret-backends",
	"add-secret-backend",
	"update-secret-backend",
	"remove-secret-backend",
	"show-secret-backend",
	"model-secret-backend",
}

type Command interface {
	// cmd.Command
	SetFlags(f *gnuflag.FlagSet)
	Init(args []string) error
	Name() string
	ToolDescription() string
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

func (c *command) getContext(ctx context.Context) (*cmd.Context, error) {
	cmdCtx, err := cmd.DefaultContext()
	if err != nil {
		return nil, err
	}
	cmdCtx.Context = ctx
	return cmdCtx, nil
}

func (c *command) getContextWithOutput(ctx context.Context) (*cmd.Context, *bytes.Buffer, *bytes.Buffer, error) {
	cmdCtx, err := cmd.DefaultContext()
	if err != nil {
		return nil, nil, nil, err
	}
	cmdCtx.Context = ctx
	
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