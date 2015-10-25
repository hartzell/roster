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

type CmdInventory struct {
	CmdDefault
	List bool
	Host string
}

func CmdInventoryFactory(ui cli.Ui) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &CmdInventory{
			CmdDefault: CmdDefault{Ui: ui},
		}, nil
	}
}

func (c *CmdInventory) Run(args []string) int {
	c.InitFlagSet()
	c.FS.BoolVar(&c.List, "list", false, "Generate a full inventory (the default behavior).")
	c.FS.StringVar(&c.Host, "host", "", "Generate a host-specific inventory for this host.")
	if err := c.FS.Parse(args); err != nil {
		if err != flag.ErrHelp {
			c.Ui.Error(fmt.Sprintf("Unable to parse arguments: %s", err))
		}
		return 1
	}

	if c.List && c.Host != "" {
		c.Ui.Error("Must specify either --list or --host, not both!")
		return 1
	}

	if c.Host != "" {
		err := c.doCmdHost(c.Host)
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

func (c *CmdInventory) doCmdHost(host string) error {
	c.Ui.Output("{}")
	return nil
}

func (c *CmdInventory) doFullInventory() error {
	state, err := fetchState(c.Dir)
	if err != nil {
		return fmt.Errorf("Unable to fetchState: %s", err)
	}

	instances, err := parseState(*state)
	if err != nil {
		return fmt.Errorf("Unable to parseState: %s", err)
	}

	funcMap := template.FuncMap{
		"groups":   Groups,
		"hostvars": HostVars,
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

func (c *CmdInventory) Synopsis() string {
	return "Generate an Ansible dynamic inventory"
}
