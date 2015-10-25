package main

//go:generate esc -o templates.go templates

import (
	"os"

	"github.com/mitchellh/cli"
)

var isDev bool

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
	c := cli.NewCLI("roster", "0.0.1")
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
