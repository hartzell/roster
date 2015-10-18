package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		ui.Error(err.Error())
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
		// default command is "inventory", with some hacks for usage/synopsis
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
