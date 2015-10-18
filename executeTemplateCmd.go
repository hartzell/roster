package main

import "github.com/mitchellh/cli"

//
// Implement the "execute-template" command

type ExecuteTemplateCommand struct {
	Ui cli.Ui
}

func (c *ExecuteTemplateCommand) Run(_ []string) int {
	c.Ui.Output("Calling ExecuteTemplateCommand.Run")
	return 0
}

func (c *ExecuteTemplateCommand) Help() string {
	return "Execute a user supplied template."
}

func (c *ExecuteTemplateCommand) Synopsis() string {
	return "Execute a user supplied template."
}
