package main

//
import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mitchellh/cli"
)

// Implement the "inventory" command

type InventoryCommand struct {
	DefaultCommand
	List bool
	Host string
}

func InventoryCommandFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &InventoryCommand{
			DefaultCommand: DefaultCommand{Ui: ui},
		}, nil
	}
}

func (c *InventoryCommand) Run(args []string) int {
	c.InitFlagSet()
	c.FS.BoolVar(&c.List, "list", false, "Generate a full inventory")
	c.FS.StringVar(&c.Host, "host", "", "The host for host-specific inventory")
	if err := c.FS.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
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
			return 1
		}
		return 0
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
	state, err := fetchState(c.Dir)
	if err != nil {
		return fmt.Errorf("Unable to fetchState: %s", err)
	}

	instances, err := parseState(*state)
	if err != nil {
		return fmt.Errorf("Unable to parseState: %s", err)
	}

	funcMap := template.FuncMap{
		"groups":   groups,
		"hostvars": hostvars,
	}

	tString, err := FSString(isDev, "/templates/dynamicInventoryTemplate")
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

func (c *InventoryCommand) Synopsis() string {
	return "(s) Generate an Ansible dynamic inventory"
}
