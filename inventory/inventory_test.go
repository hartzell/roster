package inventory

import (
	"reflect"
	"testing"
)

func TestInventory(t *testing.T) {

	i, _ := NewFromJSON(input())
	SortInventory(i)
	if !reflect.DeepEqual(i, expected()) {
		t.Error("Got did not DeepEqual expected")
	}
}

func input() string {
	return `
{

    "_meta": {
        "hostvars": {
            "10.29.92.104": {
                "one": "a",
                "two": "b"
            },
            "10.29.92.105": {
                "one": "z",
                "two": "y"
            },
            "10.29.92.120": {
                "one": "z",
                "two": "monkey"
            }
        }
    },


    "mj-master": {
        "hosts": [
            "10.29.92.104"
        ]
    },
    "alpha": {
        "hosts": [
            "10.29.92.105",
            "10.29.92.104",
            "10.29.92.120"
        ]
    },
    "gamma": {
        "hosts": [
            "10.29.92.104"
        ]
    },
    "mj-other-0": {
        "hosts": [
            "10.29.92.105"
        ]
    },
    "omega": {
        "hosts": [
            "10.29.92.120",
            "10.29.92.105"
        ]
    },
    "mj-other-1": {
        "hosts": [
            "10.29.92.120"
        ]
    }
}
`
}

func expected() *Inventory {
	return &Inventory{
		HostVars: HostVars{
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
		Groups: Groups{
			"omega": Group{
				Hosts: []string{"10.29.92.105", "10.29.92.120"},
			},
			"mj-other-1": Group{
				Hosts: []string{"10.29.92.120"},
			},
			"mj-master": Group{
				Hosts: []string{"10.29.92.104"},
			},
			"alpha": Group{
				Hosts: []string{"10.29.92.104", "10.29.92.105", "10.29.92.120"},
			},
			"gamma": Group{
				Hosts: []string{"10.29.92.104"},
			},
			"mj-other-0": Group{
				Hosts: []string{"10.29.92.105"},
			},
		},
	}
}
