package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hartzell/roster/inventory"
	"github.com/mitchellh/cli"
)

func TestMain(t *testing.T) {
	cwd, err := os.Getwd()
	err = os.Chdir("misc")
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

	if !reflect.DeepEqual(inv, expected_list_inventory()) {
		t.Error(fmt.Sprintf("Output was not as expected.\n%s", ui.OutputWriter.String()))
	}
}

func TestMainWithDir(t *testing.T) {
	ui := new(cli.MockUi)
	exitStatus, err := doIt(ui, []string{"-dir", "misc"})

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

	if !reflect.DeepEqual(inv, expected_list_inventory()) {
		t.Error(fmt.Sprintf("Output was not as expected.\n%s", ui.OutputWriter.String()))
	}
}

func expected_list_inventory() *inventory.Inventory {
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
