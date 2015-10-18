package main

import (
	"flag"

	"github.com/mitchellh/cli"
)

//
// Implement the default command (inventory, except help is different)

type DefaultInventoryCommand struct {
	List bool
	Host string
	Ui   cli.Ui
	cli  cli.CLI
}

func (c *DefaultInventoryCommand) Run(args []string) int {
	ic := InventoryCommand{
		List: c.List,
		Host: c.Host,
		Ui:   c.Ui,
	}

	// FOWL/FOUL...
	// mimic the inventory commands arg passing here so that we can do a useful
	// help message if someone gave a bogus arg to the default command.
	// TODO: refactor this so that default inventory and inventory share
	// the same bit of code

	// TODO: this still screws up `./roster --host moose`, it things
	// that moose is a subcommand instead of an arg to the default command.
	// https://github.com/mitchellh/cli/blob/master/cli.go#L165-L172

	cmdFlags := flag.NewFlagSet("defaultinventory", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.BoolVar(&ic.List, "list", false, "Generate a full inventory")
	cmdFlags.StringVar(&ic.Host, "host", "", "The host for host-specific inventory")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	return ic.Run(args)
}

func (c *DefaultInventoryCommand) Help() string {
	return c.cli.HelpFunc(c.cli.Commands) + "\n"
}

func (c *DefaultInventoryCommand) Synopsis() string {
	return "(default command is 'inventory')"
}
