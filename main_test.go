package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hartzell/roster/inventory"
	"github.com/mitchellh/cli"
)

func TestOpenStack(t *testing.T) {
	cwd, err := os.Getwd()
	err = os.Chdir("misc/fixtures/openstack")
	defer os.Chdir(cwd)

	ui := new(cli.MockUi)
	exitStatus, err := doIt(ui, []string{})

	if exitStatus != 0 {
		t.Error(fmt.Sprintf("exitStatus was %d, expected 0", exitStatus))
	}
	if err != nil {
		t.Error(fmt.Sprintf("err was \"%s\", expected nil", err))
	}

	inv, err := inventory.NewFromJSON(ui.OutputWriter.String())
	if err != nil {
		fmt.Println("Unable to create inventory from output", err)
	}
	inventory.SortInventory(inv)

	if !reflect.DeepEqual(inv, os_expected_list_inventory()) {
		t.Error(fmt.Sprintf("Output was not as expected.\n%s", ui.OutputWriter.String()))
	}
}

func TestOpenStackWithDir(t *testing.T) {
	ui := new(cli.MockUi)
	exitStatus, err := doIt(ui, []string{"-dir", "misc/fixtures/openstack"})

	if exitStatus != 0 {
		t.Error(fmt.Sprintf("exitStatus was %d, expected 0", exitStatus))
	}
	if err != nil {
		t.Error(fmt.Sprintf("err was \"%s\", expected nil", err))
	}

	inv, err := inventory.NewFromJSON(ui.OutputWriter.String())
	if err != nil {
		fmt.Println("Unable to create inventory from output", err)
	}
	inventory.SortInventory(inv)

	if !reflect.DeepEqual(inv, os_expected_list_inventory()) {
		t.Error(fmt.Sprintf("Output was not as expected.\n%s", ui.OutputWriter.String()))
	}
}

func os_expected_list_inventory() *inventory.Inventory {
	return &inventory.Inventory{
		HostVars: inventory.HostVars{
			"10.29.92.104": {
				"one": "a",
				"two": "b",
			},
			"10.29.92.105": {
				"one": "z",
				"two": "y",
			},
			"10.29.92.120": {
				"one": "z",
				"two": "monkey",
			},
		},
		Groups: inventory.Groups{
			"omega": inventory.Group{
				Hosts: []string{"10.29.92.105", "10.29.92.120"},
			},
			"mj-other-1": inventory.Group{
				Hosts: []string{"10.29.92.120"},
			},
			"mj-master": inventory.Group{
				Hosts: []string{"10.29.92.104"},
			},
			"alpha": inventory.Group{
				Hosts: []string{"10.29.92.104", "10.29.92.105", "10.29.92.120"},
			},
			"gamma": inventory.Group{
				Hosts: []string{"10.29.92.104"},
			},
			"mj-other-0": inventory.Group{
				Hosts: []string{"10.29.92.105"},
			},
		},
	}
}

func TestDigitalOcean(t *testing.T) {
	cwd, err := os.Getwd()
	err = os.Chdir("misc/fixtures/digitalocean")
	defer os.Chdir(cwd)

	ui := new(cli.MockUi)
	exitStatus, err := doIt(ui, []string{})

	if exitStatus != 0 {
		t.Error(fmt.Sprintf("exitStatus was %d, expected 0", exitStatus))
	}
	if err != nil {
		t.Error(fmt.Sprintf("err was \"%s\", expected nil", err))
	}

	inv, err := inventory.NewFromJSON(ui.OutputWriter.String())
	if err != nil {
		fmt.Println("Unable to create inventory from output", err)
	}
	inventory.SortInventory(inv)

	if !reflect.DeepEqual(inv, do_expected_list_inventory()) {
		t.Error(fmt.Sprintf("Output was not as expected.\n%s", ui.OutputWriter.String()))
	}
}

func do_expected_list_inventory() *inventory.Inventory {
	return &inventory.Inventory{
		Groups: inventory.Groups{
			"puppy": inventory.Group{
				Hosts: []string{"104.236.187.205"},
			},
			"kitty": inventory.Group{
				Hosts: []string{"159.203.251.124"},
			},
			"master": inventory.Group{
				Hosts: []string{"107.170.194.163"},
			},
		},
	}
}
