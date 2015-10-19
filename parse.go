package main

import (
	"errors"

	"github.com/hashicorp/terraform/terraform"
)

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
