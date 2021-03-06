package main

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

//
// Implement the "dump-template" command

type CmdDumpTemplate struct {
	CmdDefault
	Template string
}

func CmdDumpTemplateFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &CmdDumpTemplate{
			CmdDefault: CmdDefault{Ui: ui},
		}, nil
	}
}

func (c *CmdDumpTemplate) Run(args []string) int {
	c.InitFlagSet()
	c.FS.StringVar(&c.Template, "template", "",
		"The name of the template to dump: /templates/{dynamic-inventory,etc-hosts}.")
	if err := c.FS.Parse(args); err != nil {
		if err != flag.ErrHelp {
			c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
		}
		return 1
	}

	if c.Template == "" {
		c.Ui.Error("Missing template argument\n" + c.Help())
		return 1
	}

	tString, err := FSString(isDev, c.Template)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to read template: %s", err))
	}

	c.Ui.Output(tString)
	return 0
}

func (c *CmdDumpTemplate) Synopsis() string {
	return "Dump one of roster's built in templates."
}
