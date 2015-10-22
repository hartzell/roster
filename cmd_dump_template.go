package main

import (
	"flag"
	"fmt"
	"os"

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
