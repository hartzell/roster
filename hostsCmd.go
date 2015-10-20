package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mitchellh/cli"
)

//
// Implement the "hosts" command

type HostsCommand struct {
	Ui cli.Ui
}

func (c *HostsCommand) Run(_ []string) int {
	state, err := fetchState(".")
	if err != nil {
		return 1
	}

	instances, err := parseState(*state)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse state file: %s", err))
		return 1
	}

	tString, err := FSString(false, "/templates/etcHostsTemplate")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to read etcHostsTemplate: %s", err))
		return 1
	}

	t, err := template.New("etcHostsTemplate").Parse(tString)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse etcHostsTemplate: %s", err))
		return 1
	}

	output := bytes.NewBuffer([]byte{})
	err = t.Execute(output, instances)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to execute hosts template: %s", err))
		return 1
	}

	c.Ui.Output(output.String())
	return 0
}

func (c *HostsCommand) Help() string {
	return "Generate an /etc/hosts fragment for the Terraform instances"
}

func (c *HostsCommand) Synopsis() string {
	return "Generate an /etc/hosts fragment for the Terraform instances"
}
