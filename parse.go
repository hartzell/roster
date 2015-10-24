package main

import (
	"errors"
	"regexp"

	"github.com/hashicorp/terraform/builtin/providers/openstack"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type instanceInfo struct {
	Name     string
	Address  string
	Groups   []string
	HostVars map[string]string
}

func parseState(state terraform.State) ([]*instanceInfo, error) {

	instances := []*instanceInfo{}
	for _, m := range state.Modules {
		for _, rs := range m.Resources {
			switch rs.Type {
			case "openstack_compute_instance_v2":
				info, err := parse_os_compute_instance_v2(rs)
				if err != nil {
					return []*instanceInfo{}, errors.New("Unable to parse" + "SHITE")
				}
				instances = append(instances, info)
			}
		}
	}
	return instances, nil
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
