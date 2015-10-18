package main

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

//
// Implement the "dump-template" command

type DumpTemplateCommand struct {
	Template string
	Ui       cli.Ui
}

func (c *DumpTemplateCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("dumpTemplate", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Template, "template", "", "The name of the template to dump.")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if c.Template == "" {
		c.Ui.Output(c.Help())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Calling DumpTemplateCommand.Run: %s", c.Template))
	return 0
}

func (c *DumpTemplateCommand) Help() string {
	return "Dump one of roster's built in templates."
}

func (c *DumpTemplateCommand) Synopsis() string {
	return "Dump one of roster's built in templates."
}
