package main

import (
	"fmt"
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Node struct {
	inner *cfg.Node
}

func (n *Node) Label() int {
	return n.inner.Label
}

func (n *Node) PredsNotTaken() []int {
	preds := make([]int, len(n.inner.PredsNotTaken))
	for i, p := range n.inner.PredsNotTaken {
		preds[i] = p.Label
	}
	return preds
}

func (n *Node) PredsTaken() []int {
	preds := make([]int, len(n.inner.PredsTaken))
	for i, p := range n.inner.PredsTaken {
		preds[i] = p.Label
	}
	return preds
}

func (n *Node) SuccsNotTaken() []int {
	if n.inner.SuccNotTaken == nil {
		return []int{}
	}

	return []int{n.inner.SuccNotTaken.Label}
}

func (n *Node) SuccsTaken() []int {
	if n.inner.SuccTaken == nil {
		return []int{}
	}

	return []int{n.inner.SuccTaken.Label}
}

func (n *Node) Get() dfa.Stmt {
	return n.inner.Expr
}

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

func main() {
	f, err := ioutil.ReadFile("program.stringlang")
	if err != nil {
		panic(err)
	}
	expr, err := stringlang.Parse(f)
	if err != nil {
		panic(err)
	}
	args := os.Args[1:]
	ctx := stringlang.NewContext(args, map[string]func([]string) string{})
	fmt.Println("Result:")
	res := string(expr.Eval(ctx))
	res = strings.ReplaceAll(res, `\n`, "\n")
	fmt.Println(res)

	prog := expr.(ast.Program)
	prog = optimizer.Normalize(prog)
	graph, _ := cfg.New(prog)

	ids := make([]int, 0)
	idToNode := make(map[int]dfa.Node)

	graph.Visit(func(node *cfg.Node) {
		ids = append(ids, node.Label)

		dfaNode := new(Node)
		dfaNode.inner = node
		idToNode[node.Label] = dfaNode
	})

	sort.Ints(ids)

	merge := func(am1F, am2F dfa.Fact) dfa.Fact {
		am1 := am1F.(Set)
		am2 := am2F.(Set)
		return am1.Union(am2)
	}

	// TODO transformer that sends same flows to both outs
	flow := func(setF dfa.Fact, nodeF dfa.Node) (fallThrough dfa.Fact, branchOut dfa.Fact) {
		set := setF.(Set)
		node := nodeF.(*Node)
		expr := node.inner.Expr

		gen := Set(ast.UsedVars([]ast.Expr{expr}))
		kill := make(Set)

		if val, ok := expr.(ast.Assn); ok {
			kill = setFrom(string(val.V))
		}

		res := gen.Union(set.Except(kill))

		fmt.Println()
		fmt.Println("Flowing through:", expr)
		fmt.Println("Gen:", gen)
		fmt.Println("UNION (")
		fmt.Println(set)
		fmt.Println("EXCEPT")
		fmt.Println(kill)
		fmt.Println(")")
		fmt.Println("=", res)


		return res, res
	}

	bottom := make(Set)

	in, out := dfa.RunBackward([]int{}, ids, idToNode, merge, flow, bottom, bottom)

	for _, id := range ids {
		factF := out[id]
		fact := factF.(Set)
		nodeF := idToNode[id]
		node := nodeF.(*Node)
		fmt.Println()
		fmt.Println(node.Label(), ": ", node.inner.Expr.String())
		fmt.Println(fact)

	}

	fmt.Println()
	fmt.Println("Live before first node:")
	// These are uninitialized variables
	fmt.Println(in[graph.Entry.Label].String())
}
