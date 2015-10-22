package main

//go:generate esc -o templates.go templates

import (
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	args := os.Args[1:]

	exitStatus, err := doIt(ui, args)
	if exitStatus != 0 && err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}

func doIt(ui cli.Ui, args []string) (exitStatus int, err error) {

	c := cli.NewCLI("roster", "0.0.1")

	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"inventory": func() (cli.Command, error) {
			return &InventoryCommand{
				Ui: ui,
			}, nil
		},
		// default command is "inventory", with some hacks for usage/synopsis
		"": func() (cli.Command, error) {
			return &DefaultInventoryCommand{
				Ui:  ui,
				cli: *c,
			}, nil
		},
		"hosts": func() (cli.Command, error) {
			return &HostsCommand{
				Ui: ui,
			}, nil
		},
		"dump-template": func() (cli.Command, error) {
			return &DumpTemplateCommand{
				Ui: ui,
			}, nil
		},
		"execute-template": func() (cli.Command, error) {
			return &ExecuteTemplateCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err = c.Run()
	return
}

type instanceInfo struct {
	Name     string
	Address  string
	Groups   []string
	HostVars map[string]string
}
