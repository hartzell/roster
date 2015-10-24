// Package inventory contains types and functions to manage an Ansible
// inventory (the bits that roster uses, for now).
package inventory

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mitchellh/mapstructure"
)

// HostVar is a map of variable names to variable values.
type HostVar map[string]string

// Hostvars is a map of HostVar, keyed by a hostname.
type HostVars map[string]HostVar

// A group is struct that contains a list of hosts.
type Group struct{ Hosts []string }

// Groups is a map of groups, keyed by a group name.
type Groups map[string]Group

// An Ansible Inventory consists of host variables and groups.
type Inventory struct {
	HostVars HostVars
	Groups   Groups
}

// NewFromJSON returns a *Inventory built from an Ansible dynamic
// inventory JSON string.
func NewFromJSON(j string) (*Inventory, error) {

	inv := Inventory{}

	// Unmarshal into a simple map of interface{}
	var v map[string]interface{}
	err := json.Unmarshal([]byte(j), &v)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal expected result: %s", err)
	}

	// The only thing in "_meta" section is the hostvars bit.  Use
	// Decode to convert the map into a real structure.
	meta := map[string]HostVars{}
	mapstructure.Decode(v["_meta"], &meta)
	inv.HostVars = meta["hostvars"]

	// delete the _meta section, now the rest have a regular shape: Groups
	delete(v, "_meta")
	mapstructure.Decode(v, &inv.Groups)

	return &inv, nil
}

// SortInventory sorts any slice elements (e.g. the Hosts element of
// the Groups).
// Useful for testing comparisons using reflect.DeepEqual, possibly
// other places so put it here instead of inventory_test.go
func SortInventory(i *Inventory) {
	for _, g := range i.Groups {
		sort.Strings(g.Hosts)
	}
}
