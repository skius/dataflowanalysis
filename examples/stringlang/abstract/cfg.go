package main
//
//import (
//	"fmt"
//	"github.com/skius/stringlang/ast"
//)
//
//type Node struct {
//	label       int
//	fallThrough *Node // child if the current node doesn't branch
//	branchOut   *Node // child if the current node branches (conditional jump with condition true)
//	parents     []*Node
//	expr        ast.Expr
//}
//
//type counter struct {
//	next int
//}
//
//func NewCFG(prog *ast.Program) (*Node, map[int]*Node) {
//	code := prog.Code
//	ctr := new(counter)
//	idToNode := make(map[int]*Node)
//
//	head := buildNode(code[0], ctr)
//	idToNode[head.label] = head
//	prev := head
//	remainder := code[1:]
//	exits := newFromBlock(remainder, []*Node{prev}, ctr, false, idToNode)
//	fmt.Println("Entry is: ", head)
//	fmt.Println("Exits are:")
//	for _, exit := range exits {
//		fmt.Println(exit)
//	}
//
//	//for _, expr := range *block {
//	//
//	//	if ifelse, ok := expr.(ast.IfElse); ok {
//	//		check := buildNode(ifelse.Cond, ctr)
//	//		branch := buildNode(ifelse.Then, ctr)
//	//		fall := buildNode(ifelse.Else, ctr)
//	//		check.branchOut = branch
//	//		check.fallThrough = fall
//	//		branch.parents = append(branch.parents, check)
//	//		fall.parents = append(fall.parents, check)
//	//
//	//	}
//	//}
//	return head, idToNode
//}
//
//func newFromBlock(block ast.Block, parents []*Node, ctr *counter, branch bool, idToNode map[int]*Node) (exits []*Node) {
//	prevs := parents
//	for _, expr := range block {
//
//		if ifelse, ok := expr.(ast.IfElse); ok {
//			check := buildNode(ifelse.Cond, ctr)
//			idToNode[check.label] = check
//			check.parents = prevs
//			branchExits := newFromBlock(ifelse.Then.(ast.Block), []*Node{check}, ctr, true, idToNode)
//			fallExits := newFromBlock(ifelse.Else.(ast.Block), []*Node{check}, ctr, false, idToNode)
//			if branch {
//				for i := range prevs {
//					prevs[i].branchOut = check
//				}
//				branch = false
//			} else {
//				for i := range prevs {
//					prevs[i].fallThrough = check
//				}
//			}
//			prevs = append(branchExits, fallExits...)
//
//		} else if while, ok := expr.(ast.While); ok {
//			check := buildNode(while.Cond, ctr)
//			idToNode[check.label] = check
//			check.parents = prevs
//			branchExits := newFromBlock(while.Body.(ast.Block), []*Node{check}, ctr, true, idToNode)
//			if branch {
//				for i := range prevs {
//					prevs[i].branchOut = check
//				}
//				branch = false
//			} else {
//				for i := range prevs {
//					prevs[i].fallThrough = check
//				}
//			}
//			for _, n := range branchExits {
//				check.parents = append(check.parents, n)
//				n.fallThrough = check
//			}
//			prevs = []*Node{check}
//		} else {
//			n := buildNode(expr, ctr)
//			idToNode[n.label] = n
//			n.parents = prevs
//			if branch {
//				for i := range prevs {
//					prevs[i].branchOut = n
//				}
//				branch = false
//			} else {
//				for i := range prevs {
//					prevs[i].fallThrough = n
//				}
//			}
//			prevs = []*Node{n}
//		}
//	}
//	return prevs
//}
//
//func buildNode(expr ast.Expr, ctr *counter) *Node {
//	n := new(Node)
//	n.label = ctr.incrementAndGet()
//	n.expr = expr
//	return n
//}
//
//func (ctr *counter) incrementAndGet() int {
//	ctr.next = ctr.next + 1
//	return ctr.next
//}
