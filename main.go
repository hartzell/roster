package main

//go:generate esc -o templates.go templates

import (
	"os"

	"github.com/mitchellh/cli"
)

var isDev bool

const Version = "0.0.1"

func main() {

	isDev = os.Getenv("ROSTER_DEV") == "1"

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

func doIt(ui cli.Ui, args []string) (int, error) {

	// special case `roster --list` and `roster --host hostname`.  I
	// tried doing this by making inventory the default command, but the
	// way that cli implemented you can't do '--host hostname', it tries
	// to make 'hostname' a subcommand and fails. See:
	// https://github.com/mitchellh/cli/issues/24
	//
	if (len(args) == 1 && args[0] == "--list") ||
		(len(args) == 2 && args[0] == "--host") {
		command, err := CmdInventoryFactory(ui)()
		if err != nil {
			return 1, err
		}
		return command.Run(args), nil
	}

	c := cli.NewCLI("roster", Version)
	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"inventory":        CmdInventoryFactory(ui),
		"hosts":            CmdHostFactory(ui),
		"dump-template":    CmdDumpTemplateFactory(ui),
		"execute-template": CmdExecuteTemplateFactory(ui),
	}
	exitStatus, err := c.Run()
	return exitStatus, err
}
