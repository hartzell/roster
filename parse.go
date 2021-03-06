package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/builtin/providers/cloudstack"
	"github.com/hashicorp/terraform/builtin/providers/digitalocean"
	"github.com/hashicorp/terraform/builtin/providers/openstack"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Type InstanceInfo captures the info we want for each instance,
// name, address, a slice of groups and a map of host variables.
type InstanceInfo struct {
	Name     string
	Address  string
	Groups   []string
	HostVars map[string]string
}

// Func parseState loops over the resources in a terraform.State
// instance and returns a slice of InstanceInfo and an error.  If
// there is an error, the slice of InstanceInfo will be empty.
func parseState(state terraform.State) ([]*InstanceInfo, error) {
	instances := []*InstanceInfo{}
	for _, m := range state.Modules {
		for _, rs := range m.Resources {
			var the_parser func(*terraform.ResourceState) (*InstanceInfo, error)

			switch rs.Type {
			case "openstack_compute_instance_v2":
				the_parser = parse_os_compute_instance_v2
			case "digitalocean_droplet":
				the_parser = parse_digitalocean_droplet
			case "aws_instance":
				the_parser = parse_aws_instance
			case "cloudstack_instance":
				the_parser = parse_cloudstack_instance
			}

			if the_parser != nil {
				info, err := the_parser(rs)
				if err != nil {
					return []*InstanceInfo{},
						fmt.Errorf("Unable to parse %s resource", rs.Type)
				}
				instances = append(instances, info)
			}
		}
	}
	return instances, nil
}

func parse_digitalocean_droplet(rs *terraform.ResourceState) (*InstanceInfo, error) {
	info := InstanceInfo{}

	provider := digitalocean.Provider().(*schema.Provider)
	instanceSchema := provider.ResourcesMap["digitalocean_droplet"].Schema
	stateReader := &schema.MapFieldReader{
		Schema: instanceSchema,
		Map:    schema.BasicMapReader(rs.Primary.Attributes),
	}

	nameResult, err := stateReader.ReadField([]string{"name"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read name field: %s", err)
	}
	info.Name = nameResult.ValueOrZero(instanceSchema["name"]).(string)

	accessResult, err := stateReader.ReadField([]string{"ipv4_address"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read ipv4_address field: %s", err)
	}
	info.Address = accessResult.ValueOrZero(instanceSchema["ipv4_address"]).(string)

	return &info, nil
}

// BUG(hartzell@alerce.com) parse_aws_instance doesn't do anything
// sensible with for an instance name.  What to do?

func parse_aws_instance(rs *terraform.ResourceState) (*InstanceInfo, error) {
	info := InstanceInfo{}

	provider := aws.Provider().(*schema.Provider)
	instanceSchema := provider.ResourcesMap["aws_instance"].Schema
	stateReader := &schema.MapFieldReader{
		Schema: instanceSchema,
		Map:    schema.BasicMapReader(rs.Primary.Attributes),
	}

	// nameResult, err := stateReader.ReadField([]string{"id"})
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to read id field: %s", err)
	// }
	// info.Name = nameResult.ValueOrZero(instanceSchema["id"]).(string)
	info.Name = "moose"

	accessResult, err := stateReader.ReadField([]string{"public_ip"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read a public_ip field: %s", err)
	}
	info.Address = accessResult.ValueOrZero(instanceSchema["public_ip"]).(string)
	if info.Address == "" {
		accessResult, err := stateReader.ReadField([]string{"private_ip"})
		if err != nil {
			return nil, fmt.Errorf("Unable to read a private_ip field: %s", err)
		}
		info.Address = accessResult.ValueOrZero(instanceSchema["private_ip"]).(string)
	}

	return &info, nil
}

func parse_cloudstack_instance(rs *terraform.ResourceState) (*InstanceInfo, error) {
	info := InstanceInfo{}

	provider := cloudstack.Provider().(*schema.Provider)
	instanceSchema := provider.ResourcesMap["cloudstack_instance"].Schema
	stateReader := &schema.MapFieldReader{
		Schema: instanceSchema,
		Map:    schema.BasicMapReader(rs.Primary.Attributes),
	}

	nameResult, err := stateReader.ReadField([]string{"name"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read name field: %s", err)
	}
	info.Name = nameResult.ValueOrZero(instanceSchema["name"]).(string)

	accessResult, err := stateReader.ReadField([]string{"ipaddress"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read a ipaddress field: %s", err)
	}
	info.Address = accessResult.ValueOrZero(instanceSchema["ipaddress"]).(string)

	return &info, nil
}

// Function parse_os_compute_instance_v2 uses terraform routines to
// parse info out of a terraform.ResourceState.
//
// HEADS UP: it's use of these routines is slightly underhanded (but
// better than reverse engineering the state file format...).
//
// Thanks to @apparentlymart for this bit of code.
// See: https://github.com/hashicorp/terraform/issues/3405
func parse_os_compute_instance_v2(rs *terraform.ResourceState) (*InstanceInfo, error) {
	info := InstanceInfo{}

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
			info.HostVars, err = parseVars(value.(string))
			if err != nil {
				return nil, fmt.Errorf("Unable to parse host variables: %s", err)
			}
		}
	}

	nameResult, err := stateReader.ReadField([]string{"name"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read name field: %s", err)
	}
	info.Name = nameResult.ValueOrZero(instanceSchema["name"]).(string)

	accessResult, err := stateReader.ReadField([]string{"access_ip_v4"})
	if err != nil {
		return nil, fmt.Errorf("Unable to read access_ip_v4 field: %s", err)
	}
	info.Address = accessResult.ValueOrZero(instanceSchema["access_ip_v4"]).(string)

	return &info, nil
}

// Function parseVars converts a string like "var1 = val1, var2=val2"
// into a map[string]string{"var1": "val1", "var2": "val2}
func parseVars(s string) (map[string]string, error) {
	vars := make(map[string]string)

	if len(s) > 0 {
		name_val_pairs := splitOnComma(s)
		for _, nvp := range name_val_pairs { // each name value pair (nvp)
			v, err := splitOnEqual(nvp)
			if err != nil {
				return nil, fmt.Errorf("Unable to parseVars: %s", err)
			}
			vars[v[0]] = v[1]
		}
	}
	return vars, nil
}

// Convert a string like "a, b, something" into
// []string{"a", "b", "something"}.
func splitOnComma(s string) []string {
	comma_sep := regexp.MustCompile("\\s*,\\s*")
	return comma_sep.Split(s, -1)
}

// Convert a string like "a=bb" into []string{"a", "b"}.
func splitOnEqual(s string) ([]string, error) {
	equal_sep := regexp.MustCompile("\\s*=\\s*")
	parts := equal_sep.Split(s, -1)
	if len(parts) != 2 {
		return nil,
			fmt.Errorf(
				"Unable to split \"%s\" on an equal sign and get sensible result", s)
	}
	return parts, nil
}
