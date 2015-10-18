package main

import "github.com/mitchellh/cli"

//
// Implement the "hosts" command

type HostsCommand struct {
	Ui cli.Ui
}

func (c *HostsCommand) Run(_ []string) int {
	c.Ui.Output("Calling HostsCommand.Run")
	return 0
}

func (c *HostsCommand) Help() string {
	return "Generate an Ansible dynamic inventory for a specific host (no op)."
}

func (c *HostsCommand) Synopsis() string {
	return "Generate an Ansible dynamic inventory for a specific host (no op)"
}
