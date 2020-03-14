package gics

// Group represents an industry group.
type Group uint16

// Valid validates the `Group` value.
func (g Group) Valid() bool {
	_, ok := groupKeys[g]
	return ok
}

// String satisfies `fmt.Stringer`.
func (g Group) String() string {
	if s, ok := groupKeys[g]; ok {
		return s.name
	}
	return invalid
}

// Sector returns the `Sector` this `Group` belongs to.
func (g Group) Sector() Sector {
	return Sector(g / 100)
}

// Industries returns all `Industry`s belonging to this `Group`.
func (g Group) Industries() IndustrySet {
	is := make(IndustrySet)
	gd := groupKeys[g]
	for _, i := range industries[gd.industriesStart:gd.industriesEnd] {
		is.Add(i)
	}
	return is
}

// SubIndustries returns all `SubIndustry`s belonging to this `Group`.
func (g Group) SubIndustries() SubIndustrySet {
	sis := make(SubIndustrySet)
	gd := groupKeys[g]
	for _, si := range subIndustries[gd.subIndustriesStart:gd.subIndustriesEnd] {
		sis.Add(si)
	}
	return sis
}

// ParseGroup returns the industry `Group` corresponding to a string.
// Returns an invalid `Group` if no such `Group` exists.
func ParseGroup(s string) Group {
	return groupValues[s]
}

// Groups returns a `GroupSet` containing all `Group`s.
func Groups() GroupSet {
	gs := make(GroupSet)
	for _, g := range groups {
		gs.Add(g)
	}
	return gs
}
