package checker

import (
	"bisim-checker/lts"
	"math"
)
import . "dot-parser/pair"

type requestGenerator struct {
	eqClass           string
	numberOfProcesses int
	blockSize         int
	level             int
	current           int

	//constants
	parDegree    int
	eqClasses    []string
	bisimulation lts.Bisimulation
	graphs       Pair[*Graph, *Graph]
}

func newGenerator(parDegree int, eqClasses []string, bisimulation lts.Bisimulation, graphs Pair[*Graph, *Graph]) requestGenerator {
	gen := requestGenerator{
		parDegree:    parDegree,
		bisimulation: bisimulation,
		eqClasses:    eqClasses,
		graphs:       graphs,
	}

	gen.eqClass, gen.eqClasses = gen.eqClasses[0], gen.eqClasses[1:]
	gen.level = 1
	gen.current = 0
	class := gen.bisimulation.EquivalenceClasses[gen.eqClass]
	gen.numberOfProcesses = len(class.FirstGraphProc) + len(class.SecondGraphProc)
	gen.blockSize = int(math.Ceil(float64(gen.numberOfProcesses) / float64(gen.parDegree)))

	return gen
}

func (gen *requestGenerator) generate() (InputRequest, bool) {
	input := InputRequest{}

	extracted := 0
	class := gen.bisimulation.EquivalenceClasses[gen.eqClass]
	for extracted < gen.blockSize && gen.level*gen.current < gen.numberOfProcesses {
		current := gen.level * gen.current
		if current < len(class.FirstGraphProc) {
			input.firstGraphProcesses = append(input.firstGraphProcesses, class.FirstGraphProc[current])
		} else {
			input.secondGraphProcesses = append(input.secondGraphProcesses, class.SecondGraphProc[current-len(class.FirstGraphProc)])
		}

		extracted += 1
		gen.current += 1
	}

	if extracted == 0 {
		if gen.level*gen.current >= gen.numberOfProcesses {
			gen.level *= gen.blockSize
			gen.current = 0

			if gen.level >= gen.numberOfProcesses {
				if len(gen.eqClasses) == 0 {
					return InputRequest{}, false
				}

				gen.eqClass, gen.eqClasses = gen.eqClasses[0], gen.eqClasses[1:]

				//reset variables for new equivalence class
				gen.level = 1
				gen.current = 0
				class := gen.bisimulation.EquivalenceClasses[gen.eqClass]
				gen.numberOfProcesses = len(class.FirstGraphProc) + len(class.SecondGraphProc)
				gen.blockSize = int(math.Ceil(float64(gen.numberOfProcesses) / float64(gen.parDegree)))
			}
		}

		return gen.generate()
	} else {
		return input, true
	}
}

func Emit(bisimulation lts.Bisimulation, parDegree int) error {
	eqClasses := []string{}
	for key, _ := range bisimulation.EquivalenceClasses {
		eqClasses = append(eqClasses, key)
	}

	graphs := generateGraphEqClass(bisimulation)

	gen := newGenerator(parDegree, eqClasses, bisimulation, graphs)

	input, output, termination := spawnWorkers(parDegree, graphs)

	request, finished := gen.generate()

	pending := 0

	for !finished {
		select {
		case input <- request:
			request, finished = gen.generate()
			pending += 1
		case success := <-output:
			pending -= 1
			if !success {
				msg := <-termination
				termination <- msg
				return msg
			}
		}
	}

	for pending != 0 {
		success := <-output
		pending -= 1
		if !success {
			msg := <-termination
			termination <- msg
			return msg
		}
	}

	return nil
}

func generateGraphEqClass(bisimulation lts.Bisimulation) Pair[*Graph, *Graph] {
	pair := NewPair(
		&Graph{
			Graph:     &bisimulation.FirstGraph,
			EqClasses: map[string]string{},
		},
		&Graph{
			Graph:     &bisimulation.SecondGraph,
			EqClasses: map[string]string{},
		})

	for _, eqClass := range bisimulation.EquivalenceClasses {
		for _, process := range eqClass.FirstGraphProc {
			pair.First.EqClasses[process] = eqClass.ID
		}

		for _, process := range eqClass.SecondGraphProc {
			pair.Second.EqClasses[process] = eqClass.ID
		}
	}

	return pair
}
