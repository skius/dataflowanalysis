package dataflowanalysis

import (
	"fmt"
)

// Node represents the dataflow CFG
type Node interface {
	Label() int
	Preds() []int   // Predecessors' ids
	BranchOut() int   // Successor's id if some conditional branch is taken
	FallThrough() int // Regular successor's id
	Get() Stmt
}

// Stmt is what's contained in a Node
type Stmt interface{}

// Fact is a dataflow fact
type Fact interface {
	Equals(Fact) bool
	String() string
}

type ForwardAnalysis struct {
}

func RunAnalysis(
	entryId int,
	idToNode map[int]Node,
	ids []int,
	merge func(Fact, Fact) Fact, // Meet operator
	flow func(Fact, Node) (Fact, Fact), // Flow function
	initialFlow Fact,
	entryFlow Fact,
) map[int]Fact {

	fmt.Println("running analysis...")
	facts := make(map[int]Fact) // the facts before each node
	incomingFacts := make(map[int]map[int]Fact) // incomingFacts[id1][p1] is the fact coming to id1 from p1
	branchOutFacts := make(map[int]Fact) // keeps track of previous facts to detect changes in flow
	fallThroughFacts := make(map[int]Fact) // keeps track of previous facts to detect changes in flow

	for _, id := range ids {
		facts[id] = initialFlow
		incomingFacts[id] = make(map[int]Fact)
		branchOutFacts[id] = initialFlow
		fallThroughFacts[id] = initialFlow

	}



	worklist := make(map[int]interface{}, len(ids))
	for _, id := range ids {
		worklist[id] = struct {}{}
	}

	fmt.Println("length of worklist:", len(worklist))

	for len(worklist) > 0 {
		var nodeToProcessId int
		for k := range worklist {
			nodeToProcessId = k
			break
		}
		delete(worklist, nodeToProcessId)
		nodeToProcess := idToNode[nodeToProcessId]

		fmt.Println("Processing label = ", nodeToProcessId)
		fmt.Println("incoming facts:")
		for _, f := range incomingFacts[nodeToProcessId] {
			fmt.Println(f.String())
		}

		var fact Fact
		if nodeToProcessId == entryId {
			fact = entryFlow
		} else {
			ins := make([]Fact, 0, len(incomingFacts[nodeToProcessId]))
			for _, v := range incomingFacts[nodeToProcessId] {
				ins = append(ins, v)
			}
			fact = mergeAll(merge, ins, initialFlow)
		}
		facts[nodeToProcessId] = fact

		fallThrough, branchOut := flow(fact, nodeToProcess)
		fallThroughId := nodeToProcess.FallThrough()
		branchOutId := nodeToProcess.BranchOut()

		if fallThroughId >= 0 {
			// Has FallThrough, else we can just ignore
			prevFallThrough := fallThroughFacts[nodeToProcessId]

			if !fallThrough.Equals(prevFallThrough) {
				// Flow to FallThrough changed, add it to worklist
				fmt.Println()
				fmt.Println("Fallthrough of label=", nodeToProcessId, " was: ", prevFallThrough.String())
				fmt.Println("is now: ", fallThrough.String())
				worklist[fallThroughId] = struct {}{}
				fallThroughFacts[nodeToProcessId] = fallThrough
				incomingFacts[fallThroughId][nodeToProcessId] = fallThrough
			}
		}

		if branchOutId >= 0 {
			// Has BranchOut, else we can just ignore
			prevBranchOut := branchOutFacts[nodeToProcessId]

			if !branchOut.Equals(prevBranchOut) {
				// Flow to BranchOut changed, add it to worklist
				worklist[branchOutId] = struct {}{}
				branchOutFacts[nodeToProcessId] = branchOut
				incomingFacts[branchOutId][nodeToProcessId] = branchOut
			}
		}


		//preds := nodeToProcess.Parents()
		//inFacts := make([]Fact, 0, len(preds))
		//for _, pred := range preds {
		//	inFacts = append(inFacts, facts[pred])
		//}
		//
		//inFact := mergeAll(merge, inFacts, initialFlow)
		//fallThrough, branchOut := flow(inFact, nodeToProcess)
		//prevFallThrough := facts[nodeToProcess.FallThrough()]
		//prevBranchOut := facts[nodeToProcess.BranchOut()]
		//if !branchOut.Equals(prevBranchOut) {
		//	// Flow to BranchOut changed, add it to worklist
		//	worklist = append(worklist, nodeToProcess.BranchOut())
		//}
		//if !fallThrough.Equals(prevFallThrough) {
		//	// Flow to FallThrough changed, add it to worklist
		//	worklist = append(worklist, nodeToProcess.FallThrough())
		//}
		//facts[nodeToProcessId] = inFact
	}
	return facts
}

func mergeAll(merge func(Fact, Fact) Fact, facts []Fact, initial Fact) Fact {
	if len(facts) == 0 {
		return initial
	}

	fact := facts[0]
	for _, f := range facts[1:] {
		fact = merge(fact, f)
	}
	return fact
}
