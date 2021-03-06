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

type CmdExecuteTemplate struct {
	CmdDefault
	Template string
}

func CmdExecuteTemplateFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &CmdExecuteTemplate{
			CmdDefault: CmdDefault{Ui: ui},
		}, nil
	}
}

func (c *CmdExecuteTemplate) Run(args []string) int {
	c.InitFlagSet()
	c.FS.StringVar(&c.Template, "template", "", "The filename of the template to dump.")
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

	state, err := fetchState(c.Dir)
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
		"groups":   Groups,
		"hostvars": HostVars,
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

func (c *CmdExecuteTemplate) Help() string {
	return "Execute a user supplied template."
}

func (c *CmdExecuteTemplate) Synopsis() string {
	return "Execute a user supplied template."
}
