package gics

// SubIndustry represents a sub-industry.
type SubIndustry uint32

// Valid validates the `SubIndustry` value.
func (si SubIndustry) Valid() bool {
	_, ok := subIndustryKeys[si]
	return ok
}

// String satisfies `fmt.Stringer`.
func (si SubIndustry) String() string {
	if s, ok := subIndustryKeys[si]; ok {
		return s.name
	}
	return invalid
}

// Sector returns the `Sector` this `SubIndustry` belongs to.
func (si SubIndustry) Sector() Sector {
	return si.Group().Sector()
}

// Group returns the `Group` this `SubIndustry` belongs to.
func (si SubIndustry) Group() Group {
	return si.Industry().Group()
}

// Industry returns the `Industry` this `SubIndustry` belongs to.
func (si SubIndustry) Industry() Industry {
	return Industry(si / 100)
}

// ParseSubIndustry returns the `SubIndustry` corresponding to a string.
// Returns an invalid `SubIndustry` if no such `SubIndustry` exists.
func ParseSubIndustry(s string) SubIndustry {
	return subIndustryValues[s]
}

// SubIndustries returns a `SubIndustrySet` containing all `SubIndustry`s.
func SubIndustries() SubIndustrySet {
	sis := make(SubIndustrySet)
	for _, si := range subIndustries {
		sis.Add(si)
	}
	return sis
}
