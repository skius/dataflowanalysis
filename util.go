package dataflowanalysis

// A piToPSWrapper reverses flow through path-insensitive nodes, implementing a path-sensitive interface
type piToPSWrapper struct {
	actualNode NodePI
}

func (n *piToPSWrapper) Label() int {
	return n.actualNode.Label()
}

func (n *piToPSWrapper) PredsNotTaken() []int {
	return n.actualNode.Preds()
}

func (n *piToPSWrapper) PredsTaken() []int {
	return []int{}
}

func (n *piToPSWrapper) SuccsNotTaken() []int {
	return n.actualNode.Succs()
}

func (n *piToPSWrapper) SuccsTaken() []int {
	return []int{}
}

func (n *piToPSWrapper) Get() Stmt {
	return n.actualNode.Get()
}

// A revToFwdWrapper just reverses flow through the node
type revToFwdWrapper struct {
	actualNode Node
}

func (n *revToFwdWrapper) Label() int {
	return n.actualNode.Label()
}

func (n *revToFwdWrapper) PredsNotTaken() []int {
	return n.actualNode.SuccsNotTaken()
}

func (n *revToFwdWrapper) PredsTaken() []int {
	return n.actualNode.SuccsTaken()
}

func (n *revToFwdWrapper) SuccsNotTaken() []int {
	return n.actualNode.PredsNotTaken()
}

func (n *revToFwdWrapper) SuccsTaken() []int {
	return n.actualNode.PredsTaken()
}

func (n *revToFwdWrapper) Get() Stmt {
	return n.actualNode.Get()
}
