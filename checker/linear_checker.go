package checker

import (
	"bisim-checker/lts"
	. "dot-parser/pair"
	"errors"
	"strconv"
)

type Graph struct {
	*lts.Graph
	EqClasses map[string]string
}

func ExecuteLinear(processes Pair[[]string, []string], graphs Pair[*Graph, *Graph]) error {
	if len(processes.First) == 0 && len(processes.Second) == 0 {
		return nil
	}

	if len(processes.First) == 0 {
		processes.First, processes.Second = processes.Second, processes.First
		graphs.First, graphs.Second = graphs.Second, graphs.First
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
	var firstActions []lts.Action
	var secondActions []lts.Action

	firstProcessName, firstGraph := firstProcess.Get()
	secondProcessName, secondGraph := secondProcess.Get()
	firstProcessActions := map[Pair[string, bool]]map[string]struct{}{}
	secondProcessActions := map[Pair[string, bool]]map[string]struct{}{}

	if process, present := firstGraph.Processes[firstProcessName]; !present {
		return errors.New("Process " + firstProcessName + " not found in graph " + firstGraph.ID)
	} else {
		firstActions = process.Actions
	}

	if process, present := secondGraph.Processes[secondProcessName]; !present {
		return errors.New("Process " + secondProcessName + " not found in graph " + secondGraph.ID)
	} else {
		secondActions = process.Actions
	}

	for _, action := range firstActions {
		if eqClass, present := firstGraph.EqClasses[action.TargetProcessID]; !present {
			return errors.New("Target equivalence class of process " + action.TargetProcessID + " not found in graph " + firstGraph.ID)
		} else {
			action := NewPair(action.ChannelID, action.Send)
			if eqClasses, present := firstProcessActions[action]; !present {
				firstProcessActions[action] = map[string]struct{}{eqClass: {}}
			} else {
				eqClasses[eqClass] = struct{}{}
				firstProcessActions[action] = eqClasses
			}
		}
	}

	for _, action := range secondActions {
		if secondEqClass, present := secondGraph.EqClasses[action.TargetProcessID]; !present {
			return errors.New("Target equivalence class of process " + action.TargetProcessID + " not found in graph " + secondGraph.ID)
		} else {
			if firstEqClasses, present := firstProcessActions[NewPair(action.ChannelID, action.Send)]; !present {
				return errors.New("Action send: " + strconv.FormatBool(action.Send) + " on channel " + action.ChannelID + " not found in process " + firstProcessName)
			} else if _, present := firstEqClasses[secondEqClass]; !present {
				return errors.New("Equivalence classes of the target of the action send: " +
					strconv.FormatBool(action.Send) + " on channel " +
					action.ChannelID +
					" are different in process " +
					firstProcessName)
			}
		}
	}

	for _, action := range secondActions {
		if eqClass, present := secondGraph.EqClasses[action.TargetProcessID]; !present {
			return errors.New("Target equivalence class of process " + action.TargetProcessID + " not found in graph " + secondGraph.ID)
		} else {
			action := NewPair(action.ChannelID, action.Send)
			if eqClasses, present := secondProcessActions[action]; !present {
				secondProcessActions[action] = map[string]struct{}{eqClass: {}}
			} else {
				eqClasses[eqClass] = struct{}{}
				secondProcessActions[action] = eqClasses
			}
		}
	}

	for _, action := range firstActions {
		if firstEqClass, present := firstGraph.EqClasses[action.TargetProcessID]; !present {
			return errors.New("Target equivalence class of process " + action.TargetProcessID + " not found in graph " + firstGraph.ID)
		} else {
			if secondEqClasses, present := secondProcessActions[NewPair(action.ChannelID, action.Send)]; !present {
				return errors.New("Action send: " + strconv.FormatBool(action.Send) + " on channel " + action.ChannelID + " not found in process " + secondProcessName)
			} else if _, present := secondEqClasses[firstEqClass]; !present {
				return errors.New("Equivalence classes of the target of the action send: " +
					strconv.FormatBool(action.Send) + " on channel " +
					action.ChannelID +
					" are different in process " +
					secondProcessName)
			}
		}
	}

	return nil
}
