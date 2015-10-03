package main

// https://github.com/CiscoCloud/terraform.py

// An inventory is a map of strings to things.  Most of the strings
// are group names and the things are in turn maps with keys of
// "hosts" and "vars".  **But**, one of the strings is the magic
// "_meta" and it's map with one entry, "hostvars".  Using interface{}
// is a bit too loose, but...
type inventory map[string]interface{}

// SetHostVar sets a variable to a value for a host, they end up in
// this section of the inventory:
//
//  "_meta": {
//         "hostvars": {
//             "10.29.92.104": {
//                 "one": "a",
//                 "two": "b"
//             }
//         }
//     },
func (i *inventory) SetHostVar(hostName string, varName string, varValue string) {
	if (*i)["_meta"] == nil {
		(*i)["_meta"] = make(map[string]interface{})
	}
	m := (*i)["_meta"].(map[string]interface{})

	if m["hostvars"] == nil {
		m["hostvars"] = make(map[string]map[string]string)
	}
	h := m["hostvars"].(map[string]map[string]string)

	if h[hostName] == nil {
		h[hostName] = make(map[string]string)
	}

	h[hostName][varName] = varValue
}

// SetGroupVar adds a variable to a group, they end up in this section
// of the inventory:
//
// "my_group": {
//     "vars": {
//         "one": "a",
//         "two": "b"
//     },
//     //...
// }
//
// untested, needs additions to parser.go and data in example state
// file.
func (i *inventory) SetGroupVar(groupName string, varName string, varValue string) {
	if (*i)[groupName] == nil {
		(*i)[groupName] = make(map[string]interface{})
	}
	g := (*i)[groupName].(map[string]interface{})

	if g["vars"] == nil {
		g["vars"] = make(map[string]string)
	}
	gv := g["vars"].(map[string]string)

	gv[varName] = varValue
}

// AddHostToGroup adds a host to a group, they end up in this section
// of the inventory:
//
// "my_group": {
//     "hosts": [
//         "10.29.92.104"
//         "foo.example.com"
//     ],
//     //...
// }
func (i *inventory) AddHostToGroup(hostName string, groupName string) {
	if (*i)[groupName] == nil {
		(*i)[groupName] = make(map[string]interface{})
	}
	g := (*i)[groupName].(map[string]interface{})

	if g["hosts"] == nil {
		g["hosts"] = make([]string, 0)
	}
	g["hosts"] = append(g["hosts"].([]string), hostName)
}
