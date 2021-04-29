package main

import "github.com/skius/stringlang/ast"

// isTruthyVal returns
// 1 is s is truthy
// 0 if we don't know
// -1 if s is falsy
func isTruthyVal(s *AbsString) int {
	if s.IsBottom() {
		return -1
	}
	if s.IsTop() {
		return 0
	}

	if s.Constant == "" || s.Constant == "false" {
		return -1
	}
	return 1
}

func getAllVars(p ast.Program) []string {
	seen := make(map[string]bool)
	setVarSeen(p, seen)
	vars := make([]string, 0, len(seen))
	for k := range seen {
		vars = append(vars, k)
	}
	return vars
}

func setVarSeen(expr ast.Expr, seen map[string]bool) {
	switch val := expr.(type) {
	case ast.Program:
		setVarSeen(val.Code, seen)
	case ast.Block:
		for _, e := range val {
			setVarSeen(e, seen)
		}
	case ast.Assn:
		seen[string(val.V)] = true
		setVarSeen(val.E, seen)
	case ast.Var:
		seen[string(val)] = true
	case ast.And:
		setVarSeen(val.A, seen)
		setVarSeen(val.B, seen)
	case ast.Or:
		setVarSeen(val.A, seen)
		setVarSeen(val.B, seen)
	case ast.NotEquals:
		setVarSeen(val.A, seen)
		setVarSeen(val.B, seen)
	case ast.Equals:
		setVarSeen(val.A, seen)
		setVarSeen(val.B, seen)
	case ast.Concat:
		setVarSeen(val.A, seen)
		setVarSeen(val.B, seen)
	case ast.While:
		setVarSeen(val.Cond, seen)
		setVarSeen(val.Body, seen)
	case ast.IfElse:
		setVarSeen(val.Cond, seen)
		setVarSeen(val.Then, seen)
		setVarSeen(val.Else, seen)
	case ast.Call:
		for _, e := range val.Args {
			setVarSeen(e, seen)
		}
	case ast.Index:
		setVarSeen(val.Source, seen)
		setVarSeen(val.I, seen)
	}
	return
}
