package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type gicsDefinition struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	identifier         string
	Code               int
	groupsStart        int
	groupsEnd          int
	industriesStart    int
	industriesEnd      int
	subIndustriesStart int
	subIndustriesEnd   int
}

type gicsDefinitionMap map[string]gicsDefinition

type gicsdefinitionSlice []gicsDefinition

func (ds gicsdefinitionSlice) Len() int {
	return len(ds)
}

func (ds gicsdefinitionSlice) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds gicsdefinitionSlice) Less(i, j int) bool {
	return ds[i].Code < ds[j].Code
}

func main() {
	// Read the definitions from a file
	fileContents, err := ioutil.ReadFile("./gen/gics-classification.json")
	if err != nil {
		log.Fatal(err)
	}

	var (
		baseDefinitions   gicsDefinitionMap
		definitions       [4]gicsdefinitionSlice
		outputDefinitions = make(map[int]gicsDefinition)
	)

	// Unmarshal the definitions
	if err = json.Unmarshal(fileContents, &baseDefinitions); err != nil {
		log.Fatal(err)
	}

	// group and sort the definitions
	for k, v := range baseDefinitions {
		if v.Code, err = strconv.Atoi(k); err != nil {
			log.Fatal(err)
		}
		i := len(k)/2 - 1
		v.identifier = strings.ReplaceAll(strings.Map(func(r rune) rune {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '&' {
				return r
			}
			return -1
		}, v.Name), "&", "And")
		baseDefinitions[k] = v
		definitions[i] = append(definitions[i], v)
	}
	for i := 0; i < 4; i++ {
		sort.Sort(definitions[i])
	}

	var (
		currentIndustryIndex = -1
		currentGroupIndex    = -1
		currentSectorIndex   = -1

		currentIndustry,
		currentGroup,
		currentSector int

		qSubIndustries = len(definitions[3])
		qIndustries    = len(definitions[2])
		qGroups        = len(definitions[1])
	)

	for currentSubIndustryIndex, subIndustryDef := range append(definitions[3], gicsDefinition{Code: 99999999}) {
		var (
			subIndustryCode = subIndustryDef.Code
			industryCode    = subIndustryCode / 100
			groupCode       = industryCode / 100
			sectorCode      = groupCode / 100
		)

		outputDefinitions[subIndustryCode] = subIndustryDef

		if industryCode != currentIndustry {
			if currentIndustryIndex >= 0 {
				ci := outputDefinitions[currentIndustry]
				ci.industriesEnd = currentSubIndustryIndex
				outputDefinitions[currentIndustry] = ci
			}
			currentIndustry = industryCode
			currentIndustryIndex++
			ci := baseDefinitions[strconv.Itoa(currentIndustry)]
			ci.Code = currentIndustry
			ci.subIndustriesStart = currentSubIndustryIndex
			ci.subIndustriesEnd = qSubIndustries
			outputDefinitions[currentIndustry] = ci
		}
		if groupCode != currentGroup {
			if currentGroupIndex >= 0 {
				cg := outputDefinitions[currentGroup]
				cg.subIndustriesEnd = currentSubIndustryIndex
				cg.industriesEnd = currentIndustryIndex
				outputDefinitions[currentGroup] = cg
			}
			currentGroup = groupCode
			currentGroupIndex++
			cg := baseDefinitions[strconv.Itoa(currentGroup)]
			cg.Code = currentGroup
			cg.subIndustriesStart = currentSubIndustryIndex
			cg.subIndustriesEnd = qSubIndustries
			cg.industriesStart = currentIndustryIndex
			cg.industriesEnd = qIndustries
			outputDefinitions[currentGroup] = cg
		}
		if sectorCode != currentSector {
			if currentSectorIndex >= 0 {
				cs := outputDefinitions[currentSector]
				cs.subIndustriesEnd = currentSubIndustryIndex
				cs.industriesEnd = currentIndustryIndex
				cs.groupsEnd = currentGroupIndex
				outputDefinitions[currentSector] = cs
			}
			currentSector = sectorCode
			currentSectorIndex++
			cs := baseDefinitions[strconv.Itoa(currentSector)]
			cs.Code = currentSector
			cs.subIndustriesStart = currentSubIndustryIndex
			cs.subIndustriesEnd = qSubIndustries
			cs.industriesStart = currentIndustryIndex
			cs.industriesEnd = qIndustries
			cs.groupsStart = currentGroupIndex
			cs.groupsEnd = qGroups
			outputDefinitions[currentSector] = cs
		}
	}

	t, err := template.New("gen").Parse(`// Code generated DO NOT EDIT.

package gics

// industry sectors
const (
{{ .SectorConstants }}
)

// industry groups
const (
{{ .GroupConstants }}
)

// industries
const (
{{ .IndustryConstants }}
)

// sub-industries
const (
{{ .SubIndustryConstants }}
)

var sectors = [...]Sector{ {{ .Sectors }} }

var groups = [...]Group{ {{ .Groups }} }

var industries = [...]Industry{ {{ .Industries }} }

var subIndustries = [...]SubIndustry{ {{ .SubIndustries }} }

var sectorKeys = map[Sector]sector{
{{ .SectorData }}
}

var groupKeys = map[Group]group{
{{ .GroupData }}
}

var industryKeys = map[Industry]industry{
{{ .IndustryData }}
}

var subIndustryKeys = map[SubIndustry]subIndustry{
{{ .SubIndustryData }}
}
`)
	if err != nil {
		log.Fatal(err)
	}

	var (
		templateData = struct {
			Sectors, Groups, Industries, SubIndustries,
			SectorConstants, GroupConstants, IndustryConstants, SubIndustryConstants,
			SectorData, GroupData, IndustryData, SubIndustryData string
		}{
			Sectors:              joinCodes(definitions[0]),
			Groups:               joinCodes(definitions[1]),
			Industries:           joinCodes(definitions[2]),
			SubIndustries:        joinCodes(definitions[3]),
			SectorConstants:      constants(definitions[0], "Sector", "industry sector"),
			GroupConstants:       constants(definitions[1], "Group", "industry group"),
			IndustryConstants:    constants(definitions[2], "Industry", "industry"),
			SubIndustryConstants: constants(definitions[3], "SubIndustry", "sub-industry"),
			SectorData:           sectorData(definitions[0], outputDefinitions),
			GroupData:            groupData(definitions[1], outputDefinitions),
			IndustryData:         industryData(definitions[2], outputDefinitions),
			SubIndustryData:      subIndustryData(definitions[3], outputDefinitions),
		}
		src []byte
		w   = bytes.NewBuffer(src)
	)

	if err = t.Execute(w, templateData); err != nil {
		log.Fatal(err)
	}

	if src, err = format.Source(w.Bytes()); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err = f.Write(src); err != nil {
		log.Fatal(err)
	}
}

func joinCodes(cs gicsdefinitionSlice) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		ss[i] = strconv.Itoa(c.Code)
	}
	return strings.Join(ss, ", ")
}

func constants(cs gicsdefinitionSlice, prefix, comment string) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		ss[i] = fmt.Sprintf("\t// %s%s represents the %s %s.\n\t%s%s = %d", prefix, c.identifier, c.Name, comment, prefix, c.identifier, c.Code)
	}
	return strings.Join(ss, "\n")
}

func sectorData(cs gicsdefinitionSlice, defs map[int]gicsDefinition) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		d := defs[c.Code]

		ss[i] = fmt.Sprintf(`	Sector%s: {
		groupsStart: %d,
		groupsEnd: %d,
		group: group{
			industriesStart: %d,
			industriesEnd: %d,
			industry: industry{
				name: "%s",
				subIndustriesStart: %d,
				subIndustriesEnd: %d,
			},
		},
	},`, d.identifier, d.groupsStart, d.groupsEnd, d.industriesStart, d.industriesEnd, d.Name, d.subIndustriesStart, d.subIndustriesEnd)
	}
	return strings.Join(ss, "\n")
}

func groupData(cs gicsdefinitionSlice, defs map[int]gicsDefinition) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		d := defs[c.Code]

		ss[i] = fmt.Sprintf(`	Group%s: {
		industriesStart: %d,
		industriesEnd: %d,
		industry: industry{
			name: "%s",
			subIndustriesStart: %d,
			subIndustriesEnd: %d,
		},
	},`, d.identifier, d.industriesStart, d.industriesEnd, d.Name, d.subIndustriesStart, d.subIndustriesEnd)
	}
	return strings.Join(ss, "\n")
}

func industryData(cs gicsdefinitionSlice, defs map[int]gicsDefinition) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		d := defs[c.Code]

		ss[i] = fmt.Sprintf(`	Industry%s: {
		name: "%s",
		subIndustriesStart: %d,
		subIndustriesEnd: %d,
	},`, d.identifier, d.Name, d.subIndustriesStart, d.subIndustriesEnd)
	}
	return strings.Join(ss, "\n")
}

func subIndustryData(cs gicsdefinitionSlice, defs map[int]gicsDefinition) string {
	ss := make([]string, len(cs))
	for i, c := range cs {
		d := defs[c.Code]

		ss[i] = fmt.Sprintf(`	SubIndustry%s: {
		name: "%s",
		description: "%s",
	},`, d.identifier, d.Name, d.Description)
	}
	return strings.Join(ss, "\n")
}
