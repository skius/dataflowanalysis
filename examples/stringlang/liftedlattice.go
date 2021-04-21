package main

import (
	dfa "github.com/skius/dataflowanalysis"
	"strings"
)

type AbstractMap map[string]*AbsString

func (am AbstractMap) IsBottom() bool {
	return len(am) == 0
}

func (am AbstractMap) Meet(other AbstractMap) AbstractMap {
	if am.IsBottom() || other.IsBottom() {
		return am.copy()
	}

	res := make(map[string]*AbsString)
	for k, v := range other {
		res[k] =  am.Get(k).Meet(v)
	}
	for k, v := range am {
		res[k] = other.Get(k).Meet(v)
	}
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
		res[k] =  am.Get(k).Join(v)
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
		v2, found := other[k1]
		if !found {
			return false
		}
		if !v1.Equals(v2) {
			return false
		}
	}
	// other <= am
	for k1, v1 := range other {
		v2, found := am[k1]
		if !found {
			return false
		}
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
	mappings := make([]string, 0, len(am))
	for k, v := range am {
		mappings = append(mappings, k + "=" + v.String())
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
