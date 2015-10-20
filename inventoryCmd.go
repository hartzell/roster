package main

//
import (
	"bytes"
	"flag"
	"fmt"
	"text/template"

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
		err := c.doHostInventory(c.Host)
		if err != nil {
			c.Ui.Error(err.Error())
		}
		return 1
	}

	err := c.doFullInventory()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *InventoryCommand) doHostInventory(host string) error {
	c.Ui.Output("{}")
	return nil
}

func (c *InventoryCommand) doFullInventory() error {
	state, err := fetchState(".")
	if err != nil {
		return fmt.Errorf("Unable to fetchState: %s", err)
	}

	instances, err := parseState(*state)
	if err != nil {
		return fmt.Errorf("Unable to parseState: %s", err)
	}

	funcMap := template.FuncMap{
		"groups": groups,
	}

	tString, err := FSString(false, "/templates/dynamicInventoryTemplate")
	if err != nil {
		return fmt.Errorf("Unable to read dynamicInventoryTemplate: %s", err)
	}

	t, err := template.New("dynamicInventoryTemplate").Funcs(funcMap).Parse(tString)
	if err != nil {
		return fmt.Errorf("Unable to parse dynamicInventoryTemplate: %s", err)
	}

	output := bytes.NewBuffer([]byte{})
	err = t.Execute(output, instances)
	if err != nil {
		return fmt.Errorf("Unable to execute dynamicInventoryTemplate: %s", err)
	}

	c.Ui.Output(output.String())
	return nil
}

func (c *InventoryCommand) Help() string {
	return "(h) Generate an Ansible dynamic inventory."
}

func (c *InventoryCommand) Synopsis() string {
	return "(s) Generate an Ansible dynamic inventory"
}
