package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/mitchellh/cli"
)

//
// Implement the "execute-template" command

type ExecuteTemplateCommand struct {
	Template string
	Ui       cli.Ui
}

func (c *ExecuteTemplateCommand) Run(args []string) int {
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

	state, err := fetchState(".")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to fetchState: %s", err))
		return 1
	}

	instances, err := parseState(*state)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parseState: %s", err))
		return 1
	}

	funcMap := template.FuncMap{
		"groups":   groups,
		"hostvars": hostvars,
	}

	tBytes, err := ioutil.ReadFile(c.Template)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to read '%s'", c.Template))
		return 1
	}
	tString := string(tBytes)

	t, err := template.New("executeTemplate").Funcs(funcMap).Parse(tString)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse user supplied template: %s", err))
		return 1
	}

	output := bytes.NewBuffer([]byte{})
	err = t.Execute(output, instances)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to execute user supplied template: %s", err))
		return 1
	}

	c.Ui.Output(output.String())

	return 0
}

func (c *ExecuteTemplateCommand) Help() string {
	return "Execute a user supplied template."
}

func (c *ExecuteTemplateCommand) Synopsis() string {
	return "Execute a user supplied template."
}
