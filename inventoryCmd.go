package main

//
import (
	"flag"
	"os"

	"github.com/mitchellh/cli"
)

// Implement the "inventory" command

type InventoryCommand struct {
	List bool
	Host string
	Ui   cli.Ui
}

func (c *InventoryCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("inventory", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.BoolVar(&c.List, "list", false, "Generate a full inventory")
	cmdFlags.StringVar(&c.Host, "host", "", "The host for host-specific inventory")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if c.List && c.Host != "" {
		c.Ui.Error("Must specify either --list or --host, not both!")
		return 1
	}

	if c.Host != "" {
		status, _ := doHost(c.Host)
		return status
	}

	file := "terraform.tfstate"
	if flag.Arg(0) != "" {
		file = flag.Arg(0)
	}
	f, err := os.Open(file)
	if err != nil {
		return 1
	}
	defer f.Close()

	status, _ := doList(f)
	return status
}

func (c *InventoryCommand) Help() string {
	return "(h) Generate an Ansible dynamic inventory."
}

func (c *InventoryCommand) Synopsis() string {
	return "(s) Generate an Ansible dynamic inventory"
}
