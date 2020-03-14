package gics

//go:generate go run ./gen/gen.go

type sector struct {
	group
	groupsStart byte
	groupsEnd   byte
}

type group struct {
	industry
	industriesStart byte
	industriesEnd   byte
}

type industry struct {
	name               string
	subIndustriesStart byte
	subIndustriesEnd   byte
}

type subIndustry struct {
	name        string
	description string
}

const invalid = "invalid"

var (
	sectorValues      = map[string]Sector{}
	groupValues       = map[string]Group{}
	industryValues    = map[string]Industry{}
	subIndustryValues = map[string]SubIndustry{}
)

func init() {
	for k, si := range subIndustryKeys {
		subIndustryValues[si.name] = k
	}
	for k, i := range industryKeys {
		industryValues[i.name] = k
	}
	for k, g := range groupKeys {
		groupValues[g.name] = k
	}
	for k, s := range sectorKeys {
		sectorValues[s.name] = k
	}
}
