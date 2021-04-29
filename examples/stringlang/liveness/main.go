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

// A Node is the concrete type that implements the dfa.NodePI interface
type Node struct {
	inner *cfg.Node
}

func (n *Node) Label() int {
	return n.inner.Label
}

func (n *Node) Preds() []int {
	preds := make([]int, 0, len(n.inner.PredsNotTaken)+len(n.inner.PredsTaken))
	for _, p := range n.inner.PredsNotTaken {
		preds = append(preds, p.Label)
	}
	for _, p := range n.inner.PredsTaken {
		preds = append(preds, p.Label)
	}
	return preds
}

func (n *Node) Succs() []int {
	succs := make([]int, 0, 2)
	if n.inner.SuccNotTaken != nil {
		succs = append(succs, n.inner.SuccNotTaken.Label)
	}

	if n.inner.SuccTaken != nil {
		succs = append(succs, n.inner.SuccTaken.Label)
	}

	return succs
}

func (n *Node) Get() dfa.Stmt {
	return n.inner.Expr
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
	idToNode := make(map[int]dfa.NodePI)

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

	flow := func(setF dfa.Fact, nodeF dfa.NodePI) (res dfa.Fact) {
		set := setF.(Set)
		node := nodeF.(*Node)
		expr := node.inner.Expr

		gen := Set(ast.UsedVars([]ast.Expr{expr}))
		kill := make(Set)

		if val, ok := expr.(ast.Assn); ok {
			kill = setFrom(string(val.V))
		}

		res = gen.Union(set.Except(kill))

		return res
	}

	// We start out with just the empty set
	bottom := make(Set)

	in, out := dfa.RunBackwardPI(ids, idToNode, merge, flow, bottom)

	// Print computed liveness
	for _, id := range ids {
		factOutF := out[id]
		factOut := factOutF.(Set)
		factInF := in[id]
		factIn := factInF.(Set)
		nodeF := idToNode[id]
		node := nodeF.(*Node)
		fmt.Println()
		fmt.Println(factIn)
		fmt.Println(node.Label(), ": ", node.inner.Expr.String())
		fmt.Println(factOut)
	}
}
