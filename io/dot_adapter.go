package io

import (
	. "bisim-checker/lts"
	. "dot-parser/pair"
	"dot-parser/parser"
	. "dot-parser/result"
	"errors"
	"io"
	"os"
)

func CreateBisimulationFromFiles(firstGraphPath string, secondGraphPath string) Result[Bisimulation] {
	firstGraphReader := Make[io.Reader](os.Open(firstGraphPath))
	secondGraphReader := Make[io.Reader](os.Open(secondGraphPath))

	firstGraph := FlatMap(firstGraphReader, parser.ParseFile)
	secondGraph := FlatMap(secondGraphReader, parser.ParseFile)

	return FlatMap(Zip(firstGraph, secondGraph),
		func(pair Pair[parser.Graph, parser.Graph]) Result[Bisimulation] {
			return createBisimulation(pair.Get())
		},
	)
}

func createBisimulation(firstGraph parser.Graph, secondGraph parser.Graph) Result[Bisimulation] {
	var bisimulation Bisimulation

	graph1, err := createGraph(firstGraph).Get()
	if err != nil {
		return Err[Bisimulation](err)
	}

	graph2, err := createGraph(secondGraph).Get()
	if err != nil {
		return Err[Bisimulation](err)
	}

	bisimulation.FirstGraph = graph1.Second
	bisimulation.SecondGraph = graph2.Second

	for key, value := range graph1.First {
		if eqClass, present := bisimulation.EquivalenceClasses[key]; present {
			eqClass.FirstGraphProc = append(eqClass.FirstGraphProc, value...)
		} else {
			bisimulation.EquivalenceClasses[key] = Equivalence{
				ID:              key,
				FirstGraphProc:  value,
				SecondGraphProc: []string{},
			}
		}
	}

	for key, value := range graph2.First {
		if eqClass, present := bisimulation.EquivalenceClasses[key]; present {
			eqClass.SecondGraphProc = append(eqClass.SecondGraphProc, value...)
		} else {
			bisimulation.EquivalenceClasses[key] = Equivalence{
				ID:              key,
				FirstGraphProc:  []string{},
				SecondGraphProc: value,
			}
		}
	}

	return Ok(bisimulation)
}

func createGraph(dotGraph parser.Graph) Result[Pair[map[string][]string, Graph]] {
	var graph Graph
	var eqClasses map[string][]string

	if !dotGraph.IsDirect || dotGraph.IsStrict {
		return Err[Pair[map[string][]string, Graph]](errors.New("Input graph must be direct and not strict"))
	}

	for _, statement := range dotGraph.Statements {
		switch val := statement.(type) {
		case *parser.Node:
			result, err := createProcess(*val).Get()
			eqClass, process := result.Get()
			if err != nil {
				return Err[Pair[map[string][]string, Graph]](err)
			}
			if processes, present := eqClasses[eqClass]; present {
				processes = append(processes, process.ProcessID)
			} else {
				eqClasses[eqClass] = []string{process.ProcessID}
			}
			graph.Processes[process.ProcessID] = process
		}
	}

	for _, statement := range dotGraph.Statements {
		switch val := statement.(type) {
		case *parser.Edge:
			result, err := createAction(*val).Get()
			sourceProcess, action := result.Get()
			if err != nil {
				return Err[Pair[map[string][]string, Graph]](err)
			}
			if process, present := graph.Processes[sourceProcess]; !present {
				return Err[Pair[map[string][]string, Graph]](errors.New("The process in the action does not exists"))
			} else {
				process.Actions = append(process.Actions, action)
			}
		}
	}

	return Ok(NewPair(eqClasses, graph))
}

func createProcess(node parser.Node) Result[Pair[string, Process]] {
	var process Process
	eqClass := ""

	for _, attributeMap := range node.Attributes {
		if value, present := attributeMap["class"]; present {
			eqClass = value
		}
	}

	if eqClass == "" {
		return Err[Pair[string, Process]](errors.New("Equivalence class missing on node"))
	}

	process.ProcessID = node.ID.Name

	return Ok(NewPair(eqClass, process))
}

func createAction(edge parser.Edge) Result[Pair[string, Action]] {
	var action Action
	foundChannel := false
	foundSend := false

	for _, attributeMap := range edge.Attributes {
		if channel, present := attributeMap["channel"]; present {
			foundChannel = true
			action.ChannelID = channel
		}
		if send, present := attributeMap["send"]; present {
			if send == "t" {
				action.Send = true
				foundSend = true
			} else if send == "f" {
				action.Send = false
				foundSend = true
			}
		}
	}

	if !foundSend || !foundChannel {
		return Err[Pair[string, Action]](errors.New("Attributes missing on edge"))
	}

	action.TargetProcessID = edge.Rnode.Name

	return Ok(NewPair(edge.Lnode.Name, action))
}
