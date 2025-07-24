package jujuadapter

import (
	"log"

	"github.com/juju/juju/cmd/juju/status"
	"github.com/juju/juju/internal/cmd"
)

var commandList = []string{
	"status",
}

type Command interface {
	cmd.Command
	Name() string
	ToolDescription() string
}

func NewCommand(name string) Command {
	jcmd := getJujuCommand(name)
	return &command{
		Command: jcmd,
		info:    jcmd.Info(),
	}
}

func getJujuCommand(name string) cmd.Command {
	switch name {
	case "status":
		return status.NewStatusCommand()
	default:
		log.Panic("Unknown command")
	}
	return nil
}

type command struct {
	cmd.Command
	info *cmd.Info
}

func (c *command) Name() string {
	return c.info.Name
}

func (c *command) ToolDescription() string {
	return c.info.Purpose
}
