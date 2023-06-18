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

func MakeEquivalence() Equivalence {
	return Equivalence{
		ID:              "",
		FirstGraphProc:  []string{},
		SecondGraphProc: []string{},
	}
}

func MakeBisimulation() Bisimulation {
	return Bisimulation{
		FirstGraph:         MakeGraph(),
		SecondGraph:        MakeGraph(),
		EquivalenceClasses: map[string]Equivalence{},
	}
}
