package main

type Group struct {
	Name  string
	Hosts []string
}
type Groups map[string]*Group

func groups(instances []*instanceInfo, groupJustForHost bool) []*Group {
	groups := Groups{}
	for _, i := range instances {
		if groupJustForHost {
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

type variable struct {
	Name  string
	Value string
}
type HostVar struct {
	Host string
	Vars []variable
}

func hostvars(instances []*instanceInfo) []*HostVar {
	theVars := []*HostVar{}

	for _, i := range instances {
		if len(i.HostVars) != 0 {
			h := &HostVar{
				Host: i.Address,
			}
			for name, value := range i.HostVars {
				v := variable{Name: name, Value: value}
				h.Vars = append(h.Vars, v)
				// 	// groups[g] = append(groups[g], i.Address)
				// 	if groups[g] == nil {
				// 		groups[g] = &Group{Name: g}
				// 	}
				// 	groups[g].Hosts = append(groups[g].Hosts, i.Address)
			}
			theVars = append(theVars, h)
		}
	}
	// apparently it's faster to use make and index, but it's small, so...
	return theVars
}
