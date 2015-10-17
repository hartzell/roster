package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/builtin/providers/openstack"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/cli"
)

func main() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	args := os.Args[1:]

	exitStatus, err := doIt(ui, args)
	if exitStatus != 0 && err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func doIt(ui cli.Ui, args []string) (exitStatus int, err error) {

	c := cli.NewCLI("roster", "0.0.1")

	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"inventory": func() (cli.Command, error) {
			return &InventoryCommand{
				Ui: ui,
			}, nil
		},
		// default command is "inventory"
		"": func() (cli.Command, error) {
			return &DefaultInventoryCommand{
				Ui:  ui,
				cli: *c,
			}, nil
		},
		"hosts": func() (cli.Command, error) {
			return &HostsCommand{
				Ui: ui,
			}, nil
		},
		"dump-template": func() (cli.Command, error) {
			return &DumpTemplateCommand{
				Ui: ui,
			}, nil
		},
		"execute-template": func() (cli.Command, error) {
			return &ExecuteTemplateCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err = c.Run()
	return
}

//
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
	return ic.Run(args)
}

func (c *DefaultInventoryCommand) Help() string {
	return c.cli.HelpFunc(c.cli.Commands) + "\n"
}

func (c *DefaultInventoryCommand) Synopsis() string {
	return ""
}

//
// Implement the "hosts" command

type HostsCommand struct {
	Ui cli.Ui
}

func (c *HostsCommand) Run(_ []string) int {
	c.Ui.Output("Calling HostsCommand.Run")
	return 0
}

func (c *HostsCommand) Help() string {
	return "Generate an Ansible dynamic inventory for a specific host (no op)."
}

func (c *HostsCommand) Synopsis() string {
	return "Generate an Ansible dynamic inventory for a specific host (no op)"
}

//
// Implement the "dump-template" command

type DumpTemplateCommand struct {
	Ui cli.Ui
}

func (c *DumpTemplateCommand) Run(_ []string) int {
	c.Ui.Output("Calling DumpTemplateCommand.Run")
	return 0
}

func (c *DumpTemplateCommand) Help() string {
	return "Dump one of roster's built in templates."
}

func (c *DumpTemplateCommand) Synopsis() string {
	return "Dump one of roster's built in templates."
}

//
// Implement the "execute-template" command

type ExecuteTemplateCommand struct {
	Ui cli.Ui
}

func (c *ExecuteTemplateCommand) Run(_ []string) int {
	c.Ui.Output("Calling ExecuteTemplateCommand.Run")
	return 0
}

func (c *ExecuteTemplateCommand) Help() string {
	return "Execute a user supplied template."
}

func (c *ExecuteTemplateCommand) Synopsis() string {
	return "Execute a user supplied template."
}

// -----

func doHost(host string) (exitStatus int, errorMessage string) {
	fmt.Println("{}")
	return 0, ""
}

func doList(src io.Reader) (exitStatus int, errorMessage string) {
	state, err := terraform.ReadState(src)
	if err != nil {
		return 1, "Unable to read state file"
	}

	i := inventory{}
	instances := []*instanceInfo{}
	for _, m := range state.Modules {
		for _, rs := range m.Resources {
			switch rs.Type {
			case "openstack_compute_instance_v2":
				info, _ := parse_os_compute_instance_v2(rs)
				instances = append(instances, info)
				// create a weird little group just for this instance, so that
				// plays can refer to it by "name", which is known ahead of
				// time.
				i.AddHostToGroup(info.Address, info.Name)

				// add this instance to each group that was specified in it's
				// metadata.ansible_groups list
				for _, group := range info.Groups {
					i.AddHostToGroup(info.Address, group)
				}

				// add a host variable for each one specified via
				// metadata.ansible_hostvars
				for varName, varValue := range info.HostVars {
					i.SetHostVar(info.Address, varName, varValue)
				}
			}
		}
	}
	t, err := template.ParseFiles("etcHostsTemplate")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, instances)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(i)
	if err != nil {
		return 1, "unable to json.Marshal inventory"
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		return 1, "unable to write json to stdout"
	}

	_ = gatherGroups(instances)
	return 0, ""
}

// Converts a slice of instances (each of which might belong to one or
// more groups) into a map of group names to slices of members.
func gatherGroups(instances []*instanceInfo) map[string][]string {
	groups := map[string][]string{}
	for _, i := range instances {
		for _, g := range i.Groups {
			groups[g] = append(groups[g], i.Address)
		}
	}
	spew.Dump(groups)
	return groups
}

type instanceInfo struct {
	Name     string
	Address  string
	Groups   []string
	HostVars map[string]string
}

//
// Thanks to @apparentlymart for this bit of code that pulls
// values out of the state file info.  It's a bit underhanded
// and behind terraform's public interface, but until there's
// a better way....  See:
// https://github.com/hashicorp/terraform/issues/3405
func parse_os_compute_instance_v2(rs *terraform.ResourceState) (*instanceInfo, error) {
	info := instanceInfo{}

	provider := openstack.Provider().(*schema.Provider)
	instanceSchema := provider.ResourcesMap["openstack_compute_instance_v2"].Schema
	stateReader := &schema.MapFieldReader{
		Schema: instanceSchema,
		Map:    schema.BasicMapReader(rs.Primary.Attributes),
	}

	metadataResult, err := stateReader.ReadField([]string{"metadata"})
	if err != nil {
		return nil, errors.New("Unable to read metadata from ResourceState")
	}
	m := metadataResult.ValueOrZero(instanceSchema["metadata"])
	for key, value := range m.(map[string]interface{}) {
		if key == "ansible_groups" {
			groups := splitOnComma(value.(string))
			info.Groups = append(info.Groups, groups...)
		} else if key == "ansible_hostvars" {
			info.HostVars = parseVars(value.(string))
		}
	}

	nameResult, err := stateReader.ReadField([]string{"name"})
	if err != nil {
		return nil, errors.New("dammit #2")
	}
	info.Name = nameResult.ValueOrZero(instanceSchema["name"]).(string)

	accessResult, err := stateReader.ReadField([]string{"access_ip_v4"})
	if err != nil {
		return nil, errors.New("dammit #3")
	}
	info.Address = accessResult.ValueOrZero(instanceSchema["access_ip_v4"]).(string)

	return &info, nil
}

// TODO: Don't Panic(tm).

// Convert a string like "var1 = val1, var2=val2" into a
// map[string]string{"var1": "val1", "var2": "val2}
func parseVars(s string) map[string]string {
	vars := make(map[string]string)

	if len(s) > 0 {
		name_val_pairs := splitOnComma(s)
		for _, nvp := range name_val_pairs { // each name value pair (nvp)
			err, name, value := splitOnEqual(nvp)
			if err != nil {
				panic(err)
			}
			vars[name] = value
		}
	}

	return vars
}

// TODO: these should be consistent.

// Convert a string like "a, b, something" into
// []string{"a", "b", "something"}.
func splitOnComma(s string) []string {
	comma_sep := regexp.MustCompile("\\s*,\\s*")
	return comma_sep.Split(s, -1)
}

// Convert a string like "a=bb, something" into
// []string{"a", "b", "something"}.
func splitOnEqual(s string) (error, string, string) {
	equal_sep := regexp.MustCompile("\\s*=\\s*")
	parts := equal_sep.Split(s, -1)
	if len(parts) != 2 {
		return errors.New("Multiple equal signs seen in a single assignment statement"), "", ""
	}
	return nil, parts[0], parts[1]
}
