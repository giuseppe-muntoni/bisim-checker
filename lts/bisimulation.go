package lts

type Equivalence struct {
	ID              string
	FirstGraphProc  []string
	SecondGraphProc []string
}

type Bisimulation struct {
	FirstGraph         Graph
	SecondGraph        Graph
	EquivalenceClasses map[string]Equivalence
}
