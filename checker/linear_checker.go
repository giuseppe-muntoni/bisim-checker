package checker

import (
	"bisim-checker/lts"
	. "dot-parser/pair"
)

type Graph struct {
	Graph     *lts.Graph
	EqClasses map[string]string
}

func ExecuteLinear(processes Pair[[]string, []string], graphs Pair[*Graph, *Graph]) bool {
	firstProcess := NewPair(processes.First[0], graphs.First)
	var secondProcess Pair[string, *Graph]

	for _, process := range processes.First[1:] {
		secondProcess = NewPair(process, graphs.First)
		if !executeBinary(firstProcess, secondProcess) {
			return false
		}
		firstProcess = secondProcess
	}

	for _, process := range processes.Second {
		secondProcess = NewPair(process, graphs.Second)
		if !executeBinary(firstProcess, secondProcess) {
			return false
		}
		firstProcess = secondProcess
	}

	return true
}

func executeBinary(firstProcess Pair[string, *Graph], secondProcess Pair[string, *Graph]) bool {
	var actions map[Pair[string, bool]]string
	var firstActions []lts.Action
	var secondActions []lts.Action

	if process, present := firstProcess.Second.Graph.Processes[firstProcess.First]; !present {
		return false
	} else {
		firstActions = process.Actions
	}

	if process, present := secondProcess.Second.Graph.Processes[secondProcess.First]; !present {
		return false
	} else {
		secondActions = process.Actions
	}

	if len(firstActions) != len(secondActions) {
		return false
	}

	for _, action := range firstActions {
		if eqClass, present := firstProcess.Second.EqClasses[action.TargetProcessID]; !present {
			return false
		} else {
			actions[NewPair(action.ChannelID, action.Send)] = eqClass
		}
	}

	for _, action := range secondActions {
		if secondEqClass, present := secondProcess.Second.EqClasses[action.TargetProcessID]; !present {
			return false
		} else {
			if firstEqClass, present := actions[NewPair(action.ChannelID, action.Send)]; !present || firstEqClass != secondEqClass {
				return false
			}
		}
	}

	return true
}
