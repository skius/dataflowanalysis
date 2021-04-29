package main

import (
	"fmt"
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)
import "github.com/skius/stringlang"


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

	merge := func(am1F, am2F dfa.Fact) dfa.Fact {
		am1 := am1F.(AbstractMap)
		am2 := am2F.(AbstractMap)
		return am1.Join(am2)
	}

	vars := getAllVars(prog)


	initial := make(AbstractMap)
	//Bottom element
	//for _, variable := range vars {
	//	initial[variable] = Bottom()
	//}

	entry := make(AbstractMap)
	for _, variable := range vars {
		entry[variable] = Bottom()
	}

	flow := func(amF dfa.Fact, nodeF dfa.Node) (fallThrough, branchOut dfa.Fact) {
		am := amF.(AbstractMap).copy()
		if am.IsBottom() {
			// If this node is unreachable, its children might also be
			return am, am
		}
		node := nodeF.(*Node)
		expr := node.inner.Expr
		if len(node.SuccsTaken()) == 0 {
			// We are not in a branch
			switch val := expr.(type) {
			case ast.Assn: // The only non-branching node that changes flow is an assignment
				variable := string(val.V)
				am[variable] = transform(am, val.E)
			}
			return am, nil
		}

		// We're in a branching node

		switch val := expr.(type) {
		// TODO: Move these two and other booleans to transform (which returns "true" like stringlang)
		// then check here simply isTruthyVal(transform(am, expr))
		case ast.Equals:
			left := transform(am, val.A)
			right := transform(am, val.B)
			if left.IsConstant() && right.IsConstant() {
				if left.Constant == right.Constant {
					// This is a valid statement, branchOut always taken, fallThrough never
					return initial.copy(), am
				} else {
					// This is an invalid statement, fallThrough always taken, branchOut never
					return am, initial.copy()
				}
			}
		case ast.NotEquals:
			left := transform(am, val.A)
			right := transform(am, val.B)
			if left.IsConstant() && right.IsConstant() {
				if left.Constant != right.Constant {
					// This is a valid statement, branchOut always taken, fallThrough never
					return initial.copy(), am
				} else {
					// This is an invalid statement, fallThrough always, branchOut never
					return am, initial.copy()
				}
			}
		}
		// Otherwise no information
		return am, am
	}

	sort.Ints(ids)

	in, _, _ := dfa.RunForward(
		[]int{graph.Entry.Label},
		ids,
		idToNode,
		merge,
		flow,
		initial.copy(),
		entry.copy(),
	)

	for _, id := range ids {
		factF := in[id]
		fact := factF.(AbstractMap)
		nodeF := idToNode[id]
		node := nodeF.(*Node)
		fmt.Println()
		fmt.Println(fact)
		fmt.Println(node.Label(), ": ", node.inner.Expr.String())
	}

	//fmt.Println()
	//prog := expr.(ast.Program)
	//head := NewCFG(&prog)
	//fmt.Println(head)
	//fmt.Println()
	//fmt.Println()
	//run(head)
}

//func run(n *Node) {
//	if n == nil {
//		return
//	}
//	fmt.Println()
//	fmt.Println("Visiting node label: ", n.label)
//	fmt.Println("expr: ", n.expr.String())
//	fmt.Println("FallThrough: ", n.fallThrough)
//	fmt.Println("BranchOut: ", n.branchOut)
//	fmt.Println("Parents:")
//	for _, p := range n.parents {
//		fmt.Println("Parent: label=", p.label, " string=", p.expr.String())
//	}
//	run(n.fallThrough)
//	run(n.branchOut)
//}
