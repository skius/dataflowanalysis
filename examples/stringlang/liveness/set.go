package main

import (
	dfa "github.com/skius/dataflowanalysis"
	"sort"
	"strings"
)

type Set map[string]struct{}

func (s Set) Union(s2 Set) Set {
	u := make(Set, len(s))
	for k, v := range s {
		u[k] = v
	}

	for k, v := range s2 {
		u[k] = v
	}

	return u
}

func (s Set) Except(s2 Set) Set {
	u := make(Set, len(s))
	for k, v := range s {
		u[k] = v
	}

	for k := range s2 {
		delete(u, k)
	}

	return u
}

func (s Set) Equals(otherF dfa.Fact) bool {
	other := otherF.(Set)
	// am <= other
	for k1 := range s {
		_, ok := other[k1]
		if !ok {
			return false
		}
	}
	// other <= am
	for k1 := range other {
		_, ok := s[k1]
		if !ok {
			return false
		}
	}
	return true
}

func (s Set) String() string {
	variables := make([]string, 0, len(s))
	for k := range s {
		variables = append(variables, k)
	}

	sort.Strings(variables)

	return "{ " + strings.Join(variables, ", ") + " }"
}

func setFrom(els ...string) Set {
	s := make(Set, len(els))

	for _, el := range els {
		s[el] = struct{}{}
	}

	return s
}
