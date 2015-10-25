package main

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

//
// Implement the default command (inventory, except help is different)

type CmdDefaultInventory struct {
	CmdInventory
	cli cli.CLI
}

func CmdDefaultInventoryFactory(ui cli.Ui, c *cli.CLI) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &CmdDefaultInventory{
			CmdInventory: CmdInventory{
				CmdDefault: CmdDefault{Ui: ui},
			},
			//				CmdDefault: CmdDefault{Ui: ui},
			cli: *c,
		}, nil
	}
}

func (c *CmdDefaultInventory) Run(args []string) int {
	// FOWL/FOUL...
	// mimic the inventory commands arg passing here so that we can do a useful
	// help message if someone gave a bogus arg to the default command.
	// TODO: refactor this so that default inventory and inventory share
	// the same bit of code

	// TODO: this still screws up `./roster --host moose`, it things
	// that moose is a subcommand instead of an arg to the default command.
	// https://github.com/mitchellh/cli/blob/master/cli.go#L165-L172

	c.InitFlagSet()
	c.FS.BoolVar(&c.List, "list", false, "Generate a full inventory (the default behavior).")
	c.FS.StringVar(&c.Host, "host", "", "Generate a host-specific inventory for this host.")
	if err := c.FS.Parse(args); err != nil {
		if err != flag.ErrHelp {
			c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
		}
		return 1
	}

	ic := CmdInventory{
		CmdDefault: CmdDefault{Ui: c.Ui},
	}
	return ic.Run(args)
}

func (c *CmdDefaultInventory) Help() string {
	return c.cli.HelpFunc(c.cli.Commands) + "\n"
}

func (c *CmdDefaultInventory) Synopsis() string {
	return "(default command is 'inventory')"
}
