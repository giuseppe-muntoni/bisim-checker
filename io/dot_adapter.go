package io

import (
	. "bisim-checker/lts"
	"dot-parser/parser"
	"errors"
)

func createBisimulation(firstGraph parser.Graph, secondGraph parser.Graph) Bisimulation {
	
}

func createProcess(node parser.Node) (string, Process, error) {
	var process Process
	eqClass := ""

	for _, attributeMap := range node.Attributes {
		if value, present := attributeMap["class"]; present {
			eqClass = value
		}
	}

	if eqClass == "" {
		return "", Process{}, errors.New("Equivalence class missing on node")
	}

	process.ProcessID = node.ID.Name

	return eqClass, process, nil
}

func createAction(edge parser.Edge) (string, Action, error) {
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
		return "", Action{}, errors.New("Attributes missing on edge")
	}

	action.TargetProcessID = edge.Rnode.Name

	return edge.Lnode.Name, action, nil
}
