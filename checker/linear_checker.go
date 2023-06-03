package checker

import (
	"bisim-checker/lts"
	. "dot-parser/pair"
	"errors"
)

type Graph struct {
	Graph     *lts.Graph
	EqClasses map[string]string
}

func ExecuteLinear(processes Pair[[]string, []string], graphs Pair[*Graph, *Graph]) error {
	if len(processes.First) == 0 || len(processes.Second) == 0 {
		return errors.New("")
	}

	firstProcess := NewPair(processes.First[0], graphs.First)
	var secondProcess Pair[string, *Graph]

	for _, process := range processes.First[1:] {
		secondProcess = NewPair(process, graphs.First)
		if err := executeBinary(firstProcess, secondProcess); err != nil {
			return err
		}
		firstProcess = secondProcess
	}

	for _, process := range processes.Second {
		secondProcess = NewPair(process, graphs.Second)
		if err := executeBinary(firstProcess, secondProcess); err != nil {
			return err
		}
		firstProcess = secondProcess
	}

	return nil
}

func executeBinary(firstProcess Pair[string, *Graph], secondProcess Pair[string, *Graph]) error {
	var actions map[Pair[string, bool]]string
	var firstActions []lts.Action
	var secondActions []lts.Action

	if process, present := firstProcess.Second.Graph.Processes[firstProcess.First]; !present {
		return errors.New("")
	} else {
		firstActions = process.Actions
	}

	if process, present := secondProcess.Second.Graph.Processes[secondProcess.First]; !present {
		return errors.New("")
	} else {
		secondActions = process.Actions
	}

	if len(firstActions) != len(secondActions) {
		return errors.New("")
	}

	for _, action := range firstActions {
		if eqClass, present := firstProcess.Second.EqClasses[action.TargetProcessID]; !present {
			return errors.New("")
		} else {
			actions[NewPair(action.ChannelID, action.Send)] = eqClass
		}
	}

	for _, action := range secondActions {
		if secondEqClass, present := secondProcess.Second.EqClasses[action.TargetProcessID]; !present {
			return errors.New("")
		} else {
			if firstEqClass, present := actions[NewPair(action.ChannelID, action.Send)]; !present || firstEqClass != secondEqClass {
				return errors.New("")
			}
		}
	}

	return nil
}
