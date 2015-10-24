package main

import (
	"fmt"

	"github.com/mitchellh/cli"
)

//
// Implement the default command (inventory, except help is different)

type DefaultInventoryCommand struct {
	InventoryCommand
	cli cli.CLI
}

func DefaultInventoryCommandFactory(ui cli.Ui, c *cli.CLI) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &DefaultInventoryCommand{
			InventoryCommand: InventoryCommand{
				DefaultCommand: DefaultCommand{Ui: ui},
			},
			//				DefaultCommand: DefaultCommand{Ui: ui},
			cli: *c,
		}, nil
	}
}

func (c *DefaultInventoryCommand) Run(args []string) int {
	// FOWL/FOUL...
	// mimic the inventory commands arg passing here so that we can do a useful
	// help message if someone gave a bogus arg to the default command.
	// TODO: refactor this so that default inventory and inventory share
	// the same bit of code

	// TODO: this still screws up `./roster --host moose`, it things
	// that moose is a subcommand instead of an arg to the default command.
	// https://github.com/mitchellh/cli/blob/master/cli.go#L165-L172

	c.InitFlagSet()
	c.FS.BoolVar(&c.List, "list", false, "Generate a full inventory")
	c.FS.StringVar(&c.Host, "host", "", "The host for host-specific inventory")
	if err := c.FS.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
		return 1
	}

	ic := InventoryCommand{
		DefaultCommand: DefaultCommand{Ui: c.Ui},
	}
	return ic.Run(args)
}

func (c *DefaultInventoryCommand) Help() string {
	return c.cli.HelpFunc(c.cli.Commands) + "\n"
}

func (c *DefaultInventoryCommand) Synopsis() string {
	return "(default command is 'inventory')"
}
