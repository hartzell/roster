package main

type Group struct {
	Name  string
	Hosts []string
}
type Groups map[string]*Group

func groups(instances []*instanceInfo, foo bool) []*Group {
	groups := Groups{}
	for _, i := range instances {
		if foo {
			// add a group for each individual
			if groups[i.Name] == nil {
				groups[i.Name] = &Group{Name: i.Name}
			}
			groups[i.Name].Hosts = append(groups[i.Name].Hosts, i.Address)
		}
		// walk over the individuals group and add it to them
		for _, g := range i.Groups {
			// groups[g] = append(groups[g], i.Address)
			if groups[g] == nil {
				groups[g] = &Group{Name: g}
			}
			groups[g].Hosts = append(groups[g].Hosts, i.Address)
		}
	}

	// apparently it's faster to use make and index, but it's small, so...
	result := []*Group{}
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}
