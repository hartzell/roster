package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"text/template"
)

//
// Implement the "hosts" command

type CmdHost struct {
	CmdDefault
}

func CmdHostFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &CmdHost{
			CmdDefault: CmdDefault{Ui: ui},
		}, nil
	}
}

func (c *CmdHost) Run(args []string) int {
	c.InitFlagSet()
	if err := c.FS.Parse(args); err != nil {
		if err != flag.ErrHelp {
			c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
		}
		return 1
	}

	state, err := fetchState(c.Dir)
	if err != nil {
		return 1
	}

	instances, err := parseState(*state)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse state file: %s", err))
		return 1
	}

	tString, err := FSString(isDev, "/templates/etcHostsTemplate")
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

func (c *CmdHost) Help() string {
	return "Generate an /etc/hosts fragment for the Terraform instances"
}

func (c *CmdHost) Synopsis() string {
	return "Generate an /etc/hosts fragment for the Terraform instances"
}
