package gics

// Industry represents an industry.
type Industry uint32

// Valid validates the `Industry` value.
func (i Industry) Valid() bool {
	_, ok := industryKeys[i]
	return ok
}

// String satisfies `fmt.Stringer`.
func (i Industry) String() string {
	if s, ok := industryKeys[i]; ok {
		return s.name
	}
	return invalid
}

// Sector returns the `Sector` this industry belongs to.
func (i Industry) Sector() Sector {
	return i.Group().Sector()
}

// Group returns the `Group` this industry belongs to.
func (i Industry) Group() Group {
	return Group(i / 100)
}

// SubIndustries returns all `SubIndustry`s belonging to this `Industry`.
func (i Industry) SubIndustries() SubIndustrySet {
	sis := make(SubIndustrySet)
	id := industryKeys[i]
	for _, si := range subIndustries[id.subIndustriesStart:id.subIndustriesEnd] {
		sis.Add(si)
	}
	return sis
}

// ParseIndustry returns the `Industry` corresponding to a string.
// Returns an invalid `Industry` if no such `Industry` exists.
func ParseIndustry(s string) Industry {
	return industryValues[s]
}

// Industries returns a `IndustrySet` containing all `Industry`s.
func Industries() IndustrySet {
	is := make(IndustrySet)
	for _, i := range industries {
		is.Add(i)
	}
	return is
}
