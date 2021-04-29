package dataflowanalysis

// Node represents a path-sensitive data-flow CFG
type Node interface {
	Label() int
	PredsTaken() []int    // Predecessors that branch out to the current Node
	PredsNotTaken() []int // Predecessors that fall through to the current Node
	SuccsTaken() []int    // Successors if current Node branches out
	SuccsNotTaken() []int // Successors if current Node falls through
	Get() Stmt
}

// A NodePI is a Node used for path-insensitive data-flow analyses
type NodePI interface {
	Label() int
	Preds() []int // Predecessors
	Succs() []int // Successors
	Get() Stmt
}

// Stmt is what's contained in a Node
type Stmt interface{}

// Fact is a dataflow fact
type Fact interface {
	Equals(Fact) bool
	String() string
}

// RunBackwardPI computes a path-insensitive backward data-flow analysis
func RunBackwardPI(
	ids []int,
	idToNode map[int]NodePI,
	merge func(Fact, Fact) Fact, // Meet operator
	flow func(Fact, NodePI) Fact, // Flow function
	initialFlow Fact,
) (in, out map[int]Fact) {

	idToNodeForward := make(map[int]Node, len(idToNode))

	for k, v := range idToNode {
		ps := new(piToPSWrapper)
		ps.actualNode = v
		rev := new(revToFwdWrapper)
		rev.actualNode = ps
		idToNodeForward[k] = rev
	}

	flowWrapper := func(f Fact, n Node) (Fact, Fact) {
		actual := n.(*revToFwdWrapper).actualNode.(*piToPSWrapper).actualNode
		res := flow(f, actual)

		// No flow for Taken branches because piToPSWrapper models all branches as NotTaken
		return res, nil
	}

	// Can ignore the Taken out map because we have no Taken branches
	inForward, outForwardNT, _ := RunForward([]int{}, ids, idToNodeForward, merge, flowWrapper, initialFlow, initialFlow)

	// The in flow at each node is the out flow of the reversed data-flow and vice-versa
	return outForwardNT, inForward
}

// TODO: Need to think about how to model path-sensitive backward flows
//func RunBackward(
//	entryIds []int, // entryIds don't make much sense I think
//	ids []int,
//	idToNode map[int]Node,
//	merge func(Fact, Fact) Fact, // Meet operator
//	flow func(Fact, Node) (Fact, Fact), // Flow function
//	initialFlow Fact,
//	entryFlow Fact,
//) (in, out map[int]Fact) {
//
//
//	idToNodeReverse := make(map[int]Node, len(idToNode))
//	for k, v := range idToNode {
//		rev := new(reverseNode)
//		rev.actualNode = v
//		idToNodeReverse[k] = rev
//	}
//
//	// Need to convert our reverse nodes into the real nodes such that the flow function is not exposed to reverseNode
//	revFlow := func(f Fact, node Node) (Fact, Fact) {
//		revNode := node.(*reverseNode)
//		actualNode := revNode.actualNode
//		return flow(f, actualNode)
//	}
//
//	// ignore outTaken, because our flow is path insensitive
//	revIn, revOutNotTaken, _ := RunForward(entryIds, ids, idToNodeReverse, merge, revFlow, initialFlow, entryFlow)
//
//	return revOutNotTaken, revIn
//}

// RunForwardPI computes a path-insensitive forward data-flow analysis
func RunForwardPI(
	entryIds []int,
	ids []int,
	idToNode map[int]NodePI,
	merge func(Fact, Fact) Fact, // Meet operator
	flow func(Fact, NodePI) Fact, // Flow function
	initialFlow Fact,
	entryFlow Fact,
) (in, out map[int]Fact) {
	idToNodePS := make(map[int]Node, len(idToNode))

	for k, v := range idToNode {
		ps := new(piToPSWrapper)
		ps.actualNode = v
		idToNodePS[k] = ps
	}

	flowWrapper := func(f Fact, n Node) (Fact, Fact) {
		actual := n.(*piToPSWrapper).actualNode
		res := flow(f, actual)

		// No flow for Taken branches because piToPSWrapper models all branches as NotTaken
		return res, nil
	}

	// Can ignore the Taken out map because we have no Taken branches
	inPS, outPSNT, _ := RunForward(entryIds, ids, idToNodePS, merge, flowWrapper, initialFlow, entryFlow)

	// The in flow at each node is the out flow of the reversed data-flow and vice-versa
	return inPS, outPSNT
}

// RunForward computes a path-sensitive forward data-flow analysis
func RunForward(
	entryIds []int,
	ids []int,
	idToNode map[int]Node,
	merge func(Fact, Fact) Fact, // Meet operator
	flow func(Fact, Node) (Fact, Fact), // Flow function
	initialFlow Fact,
	entryFlow Fact,
) (in, outNotTaken, outTaken map[int]Fact) {
	// The number of nodes we are working with
	n := len(ids)

	// The in and out sets for each node
	in = make(map[int]Fact, n)
	outNotTaken = make(map[int]Fact, n)
	outTaken = make(map[int]Fact, n)

	// map instead of set to avoid adding duplicates
	worklist := make(map[int]struct{}, n)

	isEntry := make(map[int]bool)

	for _, id := range ids {
		in[id] = initialFlow
		outNotTaken[id] = initialFlow
		outTaken[id] = initialFlow

		worklist[id] = struct{}{}
	}

	for _, id := range entryIds {
		in[id] = entryFlow
		isEntry[id] = true
	}

	for len(worklist) > 0 {
		// Pop a node off the worklist
		var currNodeId int
		for k := range worklist {
			currNodeId = k
			break
		}
		delete(worklist, currNodeId)

		currNode := idToNode[currNodeId]

		inFacts := make([]Fact, 0, len(currNode.PredsNotTaken())+len(currNode.PredsTaken())+1)
		for _, pred := range currNode.PredsNotTaken() {
			inFacts = append(inFacts, outNotTaken[pred])
		}
		for _, pred := range currNode.PredsTaken() {
			inFacts = append(inFacts, outTaken[pred])
		}

		if isEntry[currNodeId] {
			inFacts = append(inFacts, entryFlow)
		}

		inFact := mergeAll(merge, inFacts, initialFlow)

		in[currNodeId] = inFact
		outNotTakenFact, outTakenFact := flow(inFact, currNode)

		if outNotTakenFact != nil && !outNotTakenFact.Equals(outNotTaken[currNodeId]) {
			// Flow changed, add successors
			for _, succ := range currNode.SuccsNotTaken() {
				worklist[succ] = struct{}{}
			}

			outNotTaken[currNodeId] = outNotTakenFact
		}

		if outTakenFact != nil && !outTakenFact.Equals(outTaken[currNodeId]) {
			// Flow changed, add successors
			for _, succ := range currNode.SuccsTaken() {
				worklist[succ] = struct{}{}
			}

			outTaken[currNodeId] = outTakenFact
		}
	}

	return in, outNotTaken, outTaken
}

//func RunAnalysis(
//	entryId int,
//	idToNode map[int]Node,
//	ids []int,
//	merge func(Fact, Fact) Fact, // Meet operator
//	flow func(Fact, Node) (Fact, Fact), // Flow function
//	initialFlow Fact,
//	entryFlow Fact,
//) map[int]Fact {
//
//	fmt.Println("running analysis...")
//	facts := make(map[int]Fact)                 // the facts before each node
//	incomingFacts := make(map[int]map[int]Fact) // incomingFacts[id1][p1] is the fact coming to id1 from p1
//	branchOutFacts := make(map[int]Fact)        // keeps track of previous facts to detect changes in flow
//	fallThroughFacts := make(map[int]Fact)      // keeps track of previous facts to detect changes in flow
//
//	for _, id := range ids {
//		facts[id] = initialFlow
//		incomingFacts[id] = make(map[int]Fact)
//		branchOutFacts[id] = initialFlow
//		fallThroughFacts[id] = initialFlow
//
//	}
//
//	worklist := make(map[int]interface{}, len(ids))
//	for _, id := range ids {
//		worklist[id] = struct{}{}
//	}
//
//	fmt.Println("length of worklist:", len(worklist))
//
//	for len(worklist) > 0 {
//		var nodeToProcessId int
//		for k := range worklist {
//			nodeToProcessId = k
//			break
//		}
//		delete(worklist, nodeToProcessId)
//		nodeToProcess := idToNode[nodeToProcessId]
//
//		fmt.Println("Processing label = ", nodeToProcessId)
//		fmt.Println("incoming facts:")
//		for _, f := range incomingFacts[nodeToProcessId] {
//			fmt.Println(f.String())
//		}
//
//		var fact Fact
//		if nodeToProcessId == entryId {
//			fact = entryFlow
//		} else {
//			ins := make([]Fact, 0, len(incomingFacts[nodeToProcessId]))
//			for _, v := range incomingFacts[nodeToProcessId] {
//				ins = append(ins, v)
//			}
//			fact = mergeAll(merge, ins, initialFlow)
//		}
//		facts[nodeToProcessId] = fact
//
//		fallThrough, branchOut := flow(fact, nodeToProcess)
//		fallThroughId := nodeToProcess.FallThrough()
//		branchOutId := nodeToProcess.BranchOut()
//
//		if fallThroughId >= 0 {
//			// Has FallThrough, else we can just ignore
//			prevFallThrough := fallThroughFacts[nodeToProcessId]
//
//			if !fallThrough.Equals(prevFallThrough) {
//				// Flow to FallThrough changed, add it to worklist
//				fmt.Println()
//				fmt.Println("Fallthrough of label=", nodeToProcessId, " was: ", prevFallThrough.String())
//				fmt.Println("is now: ", fallThrough.String())
//				worklist[fallThroughId] = struct{}{}
//				fallThroughFacts[nodeToProcessId] = fallThrough
//				incomingFacts[fallThroughId][nodeToProcessId] = fallThrough
//			}
//		}
//
//		if branchOutId >= 0 {
//			// Has BranchOut, else we can just ignore
//			prevBranchOut := branchOutFacts[nodeToProcessId]
//
//			if !branchOut.Equals(prevBranchOut) {
//				// Flow to BranchOut changed, add it to worklist
//				worklist[branchOutId] = struct{}{}
//				branchOutFacts[nodeToProcessId] = branchOut
//				incomingFacts[branchOutId][nodeToProcessId] = branchOut
//			}
//		}
//
//		//preds := nodeToProcess.Parents()
//		//inFacts := make([]Fact, 0, len(preds))
//		//for _, pred := range preds {
//		//	inFacts = append(inFacts, facts[pred])
//		//}
//		//
//		//inFact := mergeAll(merge, inFacts, initialFlow)
//		//fallThrough, branchOut := flow(inFact, nodeToProcess)
//		//prevFallThrough := facts[nodeToProcess.FallThrough()]
//		//prevBranchOut := facts[nodeToProcess.BranchOut()]
//		//if !branchOut.Equals(prevBranchOut) {
//		//	// Flow to BranchOut changed, add it to worklist
//		//	worklist = append(worklist, nodeToProcess.BranchOut())
//		//}
//		//if !fallThrough.Equals(prevFallThrough) {
//		//	// Flow to FallThrough changed, add it to worklist
//		//	worklist = append(worklist, nodeToProcess.FallThrough())
//		//}
//		//facts[nodeToProcessId] = inFact
//	}
//	return facts
//}

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
