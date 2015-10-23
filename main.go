package main

//go:generate esc -o templates.go templates

import (
	"flag"
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

type DefaultCommand struct {
	Dir string
	Ui  cli.Ui
	FS  *flag.FlagSet
}

func (dc *DefaultCommand) Help() string {
	return "shite"
}

func (dc *DefaultCommand) InitFlagSet() {
	dc.FS = flag.NewFlagSet("inventory", flag.ContinueOnError)
	dc.FS.StringVar(&dc.Dir, "dir", ".", "The path to the terraform directory")
}

func doIt(ui cli.Ui, args []string) (exitStatus int, err error) {

	c := cli.NewCLI("roster", "0.0.1")

	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"inventory": func() (cli.Command, error) {
			return &InventoryCommand{
				DefaultCommand: DefaultCommand{Ui: ui},
			}, nil
		},
		// default command is "inventory", with some hacks for usage/synopsis
		"": func() (cli.Command, error) {
			return &DefaultInventoryCommand{
				InventoryCommand: InventoryCommand{
					DefaultCommand: DefaultCommand{Ui: ui},
				},
				//				DefaultCommand: DefaultCommand{Ui: ui},
				cli: *c,
			}, nil
		},
		"hosts": func() (cli.Command, error) {
			return &HostsCommand{
				DefaultCommand: DefaultCommand{Ui: ui},
			}, nil
		},
		"dump-template": func() (cli.Command, error) {
			return &DumpTemplateCommand{
				DefaultCommand: DefaultCommand{Ui: ui},
			}, nil
		},
		"execute-template": func() (cli.Command, error) {
			return &ExecuteTemplateCommand{
				DefaultCommand: DefaultCommand{Ui: ui},
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
