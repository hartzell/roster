package main

//go:generate esc -o templates.go templates

import (
	"flag"
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

type DefaultCommand struct {
	Dir string
	Ui  cli.Ui
	FS  *flag.FlagSet
}

func (dc *DefaultCommand) Help() string {
	// stub this out.  It never seems to get called, so that fact that
	// it's shared amongst all of the Commands isn't a problem.  I'm not
	// sure *why* I'm getting lucky and need to walk through it, but for
	// now, just take it.
	return ""
}

// InitFlagSet intializes the DefaultCommand's FS element and adds
// default flags.
func (dc *DefaultCommand) InitFlagSet() {
	dc.FS = flag.NewFlagSet("inventory", flag.ContinueOnError)
	dc.FS.StringVar(&dc.Dir, "dir", ".", "The path to the terraform directory")
}

func doIt(ui cli.Ui, args []string) (int, error) {
	c := cli.NewCLI("roster", "0.0.1")
	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		// default command is "inventory", with some hacks for
		// usage/synopsis, so its Factory takes diff arg than other
		// factory...
		"":                 DefaultInventoryCommandFactory(ui, c),
		"inventory":        InventoryCommandFactory(ui),
		"hosts":            HostCommandFactory(ui),
		"dump-template":    DumpTemplateCommandFactory(ui),
		"execute-template": ExecuteTemplateCommandFactory(ui),
	}

	exitStatus, err := c.Run()
	return exitStatus, err
}
