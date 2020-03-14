package gics

// Sector represents an industry sector.
type Sector byte

// Valid validates the `Sector` value.
func (s Sector) Valid() bool {
	_, ok := sectorKeys[s]
	return ok
}

// String satisfies `fmt.Stringer`.
func (s Sector) String() string {
	if s, ok := sectorKeys[s]; ok {
		return s.name
	}
	return invalid
}

// Groups returns all `Group`s belonging to this `Sector`.
func (s Sector) Groups() GroupSet {
	gs := make(GroupSet)
	sd := sectorKeys[s]
	for _, g := range groups[sd.groupsStart:sd.groupsEnd] {
		gs.Add(g)
	}
	return gs
}

// Industries returns all `Industry`s belonging to this `Sector`.
func (s Sector) Industries() IndustrySet {
	gs := make(IndustrySet)
	sd := sectorKeys[s]
	for _, i := range industries[sd.industriesStart:sd.industriesEnd] {
		gs.Add(i)
	}
	return gs
}

// SubIndustries returns all `SubIndustry`s belonging to this `Sector`.
func (s Sector) SubIndustries() SubIndustrySet {
	gs := make(SubIndustrySet)
	sd := sectorKeys[s]
	for _, i := range subIndustries[sd.subIndustriesStart:sd.subIndustriesEnd] {
		gs.Add(i)
	}
	return gs
}

// ParseSector returns the `industry Sector` corresponding to a string.
// Returns an invalid `Sector` if no such `Sector` exists.
func ParseSector(s string) Sector {
	return sectorValues[s]
}

// Sectors returns a `SectorSet` containing all `Sector`s.
func Sectors() SectorSet {
	ss := make(SectorSet)
	for _, s := range sectors {
		ss.Add(s)
	}
	return ss
}
