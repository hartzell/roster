package main

import (
	"bytes"
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
		return 1
	}

	t, err := template.ParseFiles("etcHostsTemplate")
	if err != nil {
		return 1
	}

	output := bytes.NewBuffer([]byte{})
	err = t.Execute(output, instances)
	if err != nil {
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
