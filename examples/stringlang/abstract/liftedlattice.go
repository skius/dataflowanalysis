package main

import (
	dfa "github.com/skius/dataflowanalysis"
	"sort"
	"strings"
)

/*
	The lattice of Var -> *AbsString
	Special Bottom element (empty map) that models unreachability
	A Bottom *AbsString just models uninitializedness
	So the map (\v -> Bottom) models a reachable piece of code, where no variable is initialized
 */

type AbstractMap map[string]*AbsString

func (am AbstractMap) IsBottom() bool {
	return len(am) == 0
	//for _, v := range am {
	//	if !v.IsBottom() {
	//		return false
	//	}
	//}
	//return true
}

func (am AbstractMap) Meet(other AbstractMap) AbstractMap {
	if am.IsBottom() || other.IsBottom() {
		return am.copy()
	}

	res := make(AbstractMap)
	for k, v := range other {
		res[k] = am.Get(k).Meet(v)
	}
	for k, v := range am {
		res[k] = other.Get(k).Meet(v)
	}
	//fmt.Println()
	//fmt.Println("Doing the meet: ", am.String() + " MEET " + other.String())
	//fmt.Println("results in: ", res.String())
	//fmt.Println()
	return res
}
func (am AbstractMap) Join(other AbstractMap) AbstractMap {
	if am.IsBottom() {
		return other.copy()
	}
	if other.IsBottom() {
		return am.copy()
	}

	res := make(map[string]*AbsString)
	for k, v := range other {
		res[k] = am.Get(k).Join(v)
	}
	for k, v := range am {
		res[k] = other.Get(k).Join(v)
	}
	return res
}

func (am AbstractMap) Equals(otherF dfa.Fact) bool {
	other := otherF.(AbstractMap)
	// am <= other
	for k1, v1 := range am {
		v2 := other.Get(k1)
		if !v1.Equals(v2) {
			return false
		}
	}
	// other <= am
	for k1, v1 := range other {
		v2 := am.Get(k1)
		if !v1.Equals(v2) {
			return false
		}
	}
	return true
}

func (am AbstractMap) Get(variable string) *AbsString {
	res, ok := am[variable]
	if !ok {
		return Bottom()
	}
	return res
}

func (am AbstractMap) String() string {
	if am.IsBottom() {
		return "{ <BOTTOM> }"
	}

	variables := make([]string, 0, len(am))
	for k := range am {
		variables = append(variables, k)
	}

	sort.Strings(variables)

	mappings := make([]string, 0, len(am))
	for _, variable := range variables {
		abs := am[variable]
		mappings = append(mappings, variable+"="+abs.String())
	}
	return "{ " + strings.Join(mappings, ", ") + " }"
}

func (am AbstractMap) copy() AbstractMap {
	am2 := make(map[string]*AbsString, len(am))
	for k, v := range am {
		am2[k] = v.copy()
	}
	return am2
}
