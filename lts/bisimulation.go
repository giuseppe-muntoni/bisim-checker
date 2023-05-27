package lts

type Equivalence struct {
	firstGraphProc  []string
	secondGraphProc []string
}

type Bisimulation struct {
	firstGraph         Graph
	secondGraph        Graph
	equivalenceClasses []Equivalence
}
