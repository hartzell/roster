package main

import (
	"encoding/json"
	"html/template"
)

type Group []string
type Groups map[string][]string

func groups(instances []*instanceInfo) Groups {
	groups := Groups{}
	for _, i := range instances {
		// add a group for each individual
		groups[i.Name] = append(groups[i.Name], i.Address)

		// walk over the individuals group and add it to them
		for _, g := range i.Groups {
			groups[g] = append(groups[g], i.Address)
		}
	}
	return groups
}

func quote(strings []string) []string {
	q := []string{}
	for _, s := range strings {
		// q = append(q, "\""+s+"\"")
		s, _ := json.Marshal(template.JSEscapeString(s))
		q = append(q, string(s))
	}
	return q
}

func blah(strings []string) string {
	s, _ := json.Marshal(strings)
	return string(s)
}
