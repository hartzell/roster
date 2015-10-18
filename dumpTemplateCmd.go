package main

import "github.com/mitchellh/cli"

//
// Implement the "dump-template" command

type DumpTemplateCommand struct {
	Ui cli.Ui
}

func (c *DumpTemplateCommand) Run(_ []string) int {
	c.Ui.Output("Calling DumpTemplateCommand.Run")
	return 0
}

func (c *DumpTemplateCommand) Help() string {
	return "Dump one of roster's built in templates."
}

func (c *DumpTemplateCommand) Synopsis() string {
	return "Dump one of roster's built in templates."
}
