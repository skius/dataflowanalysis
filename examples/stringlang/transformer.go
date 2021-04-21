package main

import (
	"github.com/skius/stringlang/ast"
	"strconv"
)

func transform(am AbstractMap, expr ast.Expr) *AbsString {
	switch val := expr.(type) {
	case ast.Concat:
		rhs := transform(am, val.A)
		lhs := transform(am, val.B)
		return rhs.Concat(lhs)
	case ast.Val:
		return Constant(string(val))
	case ast.Var:
		return am.Get(string(val))
	case ast.Index:
		idx := transform(am, val.I)
		src := transform(am, val.Source)
		return src.Index(idx)
	case ast.Assn:
		panic("Assignment in transformer!")
	}
	return Top()
}

func (s *AbsString) Index(idx *AbsString) *AbsString {
	if s.IsTop() || idx.IsTop() {
		return Top()
	}



	// Both are constant
	i, err := strconv.Atoi(idx.forceEval())
	if err != nil {
		return Constant("") // StringLang semantics
	}
	src := s.forceEval()
	if i >= len(src) {
		return Constant("") // StringLang semantics
	}
	return Constant(string(src[i]))
}

func (s *AbsString) Concat(other *AbsString) *AbsString {
	if s.IsTop() || other.IsTop() {
		return Top()
	}

	a := s.forceEval()
	b := other.forceEval()
	return Constant(a + b)



	//a := "" // StringLang semantics
	//if s.IsConstant() {
	//	a = s.Constant
	//}
	//b := "" // StringLang semantics
	//if other.IsConstant() {
	//	b = other.Constant
	//}
	//// Both are constant
	//return Constant(a + b)
}

func (s *AbsString) forceEval() string {
	if s.IsTop() {
		panic("forceEval on Top! Don't know what Top is!")
	}

	if s.IsBottom() {
		return "" // StringLang semantics, uninitialized variables are ""
	}
	return s.Constant
}
