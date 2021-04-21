package main

import (
	"fmt"
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/ast"
	"io/ioutil"
	"os"
	"strings"
)
import "github.com/skius/stringlang"

func (n *Node) Label() int {
	return n.label
}

func (n *Node) Preds() []int {
	preds := make([]int, len(n.parents))
	for i, p := range n.parents {
		preds[i] = p.label
	}
	return preds
}

func (n *Node) BranchOut() int {
	if n.branchOut == nil {
		return -1
	}

	return n.branchOut.label
}

func (n *Node) FallThrough() int {
	if n.fallThrough == nil {
		return -1
	}

	return n.fallThrough.label
}

func (n *Node) Get() dfa.Stmt {
	return n.expr
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
	head, idToNode := New(&prog)

	ids := make([]int, 0, len(idToNode))
	for k := range idToNode {
		ids = append(ids, k)
	}

	merge := func(am1F, am2F dfa.Fact) dfa.Fact {
		am1 := am1F.(AbstractMap)
		am2 := am2F.(AbstractMap)
		return am1.Meet(am2)
	}

	vars := getAllVars(prog)

	idToNodeGeneral := make(map[int]dfa.Node, len(idToNode))
	for k, v := range idToNode {
		idToNodeGeneral[k] = v
	}

	initial := make(AbstractMap)
	for _, variable := range vars {
		initial[variable] = Bottom()
	}

	entry := make(AbstractMap)
	for _, variable := range vars {
		entry[variable] = Top()
	}

	flow := func(amF dfa.Fact, nodeF dfa.Node) (fallThrough, branchOut dfa.Fact) {
		am := amF.(AbstractMap)
		node := nodeF.(*Node)
		fmt.Println("Flowing through expr ", node.expr.String())
		if node.label == head.label {
			// entry flow
			fmt.Println("entry!")
			return entry, entry
		}
		return am, am
	}

	facts := dfa.RunAnalysis(
			head,
			idToNodeGeneral,
			ids,
			merge,
			flow,
			initial,
		)

	for id, factF := range facts {
		fact := factF.(AbstractMap)
		node := idToNode[id]
		fmt.Println()
		fmt.Println("Fact before node label=", node.label, " and expr=", node.expr.String())
		fmt.Println(fact)
	}

	//fmt.Println()
	//prog := expr.(ast.Program)
	//head := New(&prog)
	//fmt.Println(head)
	//fmt.Println()
	//fmt.Println()
	//run(head)
}

func run(n *Node) {
	if n == nil {
		return
	}
	fmt.Println()
	fmt.Println("Visiting node label: ", n.label)
	fmt.Println("expr: ", n.expr.String())
	fmt.Println("FallThrough: ", n.fallThrough)
	fmt.Println("BranchOut: ", n.branchOut)
	fmt.Println("Parents:")
	for _, p := range n.parents {
		fmt.Println("Parent: label=", p.label, " string=", p.expr.String())
	}
	run(n.fallThrough)
	run(n.branchOut)
}
