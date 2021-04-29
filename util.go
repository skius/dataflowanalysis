package dataflowanalysis

// A reverseNode just reverses flow through the node
type reverseNode struct {
	actualNode Node
}

func (n *reverseNode) Label() int {
	return n.actualNode.Label()
}

func (n *reverseNode) PredsTaken() []int {
	return n.actualNode.SuccsTaken()
}

func (n *reverseNode) PredsNotTaken() []int {
	return n.actualNode.SuccsNotTaken()
}

func (n *reverseNode) SuccsTaken() []int {
	return n.actualNode.PredsTaken()
}

func (n *reverseNode) SuccsNotTaken() []int {
	return n.actualNode.PredsNotTaken()
}

func (n *reverseNode) Get() Stmt {
	return n.actualNode.Get()
}
