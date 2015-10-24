package main

import (
	"fmt"
	"os"
)

//
// Implement the "dump-template" command

type DumpTemplateCommand struct {
	DefaultCommand
	Template string
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

	useLocal := os.Getenv("ROSTER_DEV") == "1"
	tString, err := FSString(useLocal, c.Template)
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
