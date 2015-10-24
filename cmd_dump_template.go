package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

//
// Implement the "dump-template" command

type DumpTemplateCommand struct {
	DefaultCommand
	Template string
}

func DumpTemplateCommandFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &InventoryCommand{
			DefaultCommand: DefaultCommand{Ui: ui},
		}, nil
	}
}

func (c *DumpTemplateCommand) Run(args []string) int {
	c.InitFlagSet()
	c.FS.StringVar(&c.Template, "template", "", "The name of the template to dump.")
	if err := c.FS.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
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

func (c *DumpTemplateCommand) Help() string {
	return "(h) Dump one of roster's built in templates."
}

func (c *DumpTemplateCommand) Synopsis() string {
	return "Dump one of roster's built in templates."
}
