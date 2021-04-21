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
}

type ForwardAnalysis struct {
}

func RunAnalysis(
	entry Node,
	idToNode map[int]Node,
	ids []int,
	merge func(Fact, Fact) Fact, // Meet operator
	flow func(Fact, Node) (Fact, Fact), // Flow function
	initialFlow Fact,
) map[int]Fact {
	fmt.Println("running analysis...")
	facts := make(map[int]Fact)
	incomingFacts := make(map[int][]Fact)
	branchOutFacts := make(map[int]Fact)
	fallThroughFacts := make(map[int]Fact)

	for _, id := range ids {
		facts[id] = initialFlow
		incomingFacts[id] = []Fact{}
		branchOutFacts[id] = initialFlow
		fallThroughFacts[id] = initialFlow

	}



	worklist := make([]int, len(ids))
	copy(worklist, ids)

	fmt.Println("length of worklist:", len(worklist))

	for len(worklist) > 0 {
		nodeToProcessId := worklist[0]
		worklist = worklist[1:]
		nodeToProcess := idToNode[nodeToProcessId]

		fallThrough, branchOut := flow(facts[nodeToProcessId], nodeToProcess)
		fallThroughId := nodeToProcess.FallThrough()
		branchOutId := nodeToProcess.BranchOut()

		if fallThroughId >= 0 {
			// Has FallThrough, else we can just ignore
			prevFallThrough := facts[nodeToProcess.FallThrough()]

			if !fallThrough.Equals(prevFallThrough) {
				// Flow to FallThrough changed, add it to worklist
				worklist = append(worklist, nodeToProcess.FallThrough())
				facts[nodeToProcess.FallThrough()] = merge(fallThrough, prevFallThrough)
			}
		}

		if branchOutId >= 0 {
			// Has BranchOut, else we can just ignore
			prevBranchOut := facts[nodeToProcess.BranchOut()]

			if !branchOut.Equals(prevBranchOut) {
				// Flow to BranchOut changed, add it to worklist
				worklist = append(worklist, nodeToProcess.BranchOut())
				facts[nodeToProcess.BranchOut()] = merge(branchOut, prevBranchOut)
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
	if len(facts) == 1 {
		return facts[0]
	}

	fact := facts[0]
	for _, f := range facts[1:] {
		fact = merge(fact, f)
	}
	return fact
}
