package main

import (
	"reflect"
	"sort"
	"testing"
)

// ByAge implements sort.Interface for []*Group based on
// the Name field.
type ByName []*Group

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

// TestGroups checks that the correct []*Group is returned for a
// substantive input.
func TestGroups(t *testing.T) {
	input := []*instanceInfo{
		{Name: "i_a", Address: "192.168.1.1", Groups: []string{"g1", "g2"}},
		{Name: "i_b", Address: "192.168.1.2", Groups: []string{"g1", "g3"}},
	}
	expect := []*Group{
		&Group{Name: "g1", Hosts: []string{"192.168.1.1", "192.168.1.2"}},
		&Group{Name: "g2", Hosts: []string{"192.168.1.1"}},
		&Group{Name: "g3", Hosts: []string{"192.168.1.2"}},
	}

	got := groups(input, false)
	sort.Sort(ByName(got))
	if !reflect.DeepEqual(got, expect) {
		t.Fail()
	}

	// if you include a group for each host too, expect this...
	expect_more := []*Group{
		&Group{Name: "g1", Hosts: []string{"192.168.1.1", "192.168.1.2"}},
		&Group{Name: "g2", Hosts: []string{"192.168.1.1"}},
		&Group{Name: "g3", Hosts: []string{"192.168.1.2"}},
		&Group{Name: "i_a", Hosts: []string{"192.168.1.1"}},
		&Group{Name: "i_b", Hosts: []string{"192.168.1.2"}},
	}

	got = groups(input, true)
	sort.Sort(ByName(got))
	if !reflect.DeepEqual(got, expect_more) {
		t.Fail()
	}
}

// TestTmptyInput checks that calling groups with
// []*InstanceInfo{} returns an empty slice of results.
func TestEmptyInput(t *testing.T) {
	got := groups([]*instanceInfo{}, true)
	if len(got) != 0 {
		t.Fail()
	}
	got = groups([]*instanceInfo{}, false)
	if len(got) != 0 {
		t.Fail()
	}
}
