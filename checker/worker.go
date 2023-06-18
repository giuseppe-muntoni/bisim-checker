package checker

import . "dot-parser/pair"

type InputRequest struct {
	firstGraphProcesses  []string
	secondGraphProcesses []string
}

func work(input <-chan InputRequest, output chan<- bool, termination chan error, graphs Pair[*Graph, *Graph]) {
	for {
		select {
		case request := <-input:
			if err := ExecuteLinear(NewPair(request.firstGraphProcesses, request.secondGraphProcesses), graphs); err == nil {
				output <- true
			} else {
				termination <- err
				output <- false
			}
		case msg := <-termination:
			termination <- msg
			return
		}
	}
}

func spawnWorkers(workers int, graphs Pair[*Graph, *Graph]) (chan<- InputRequest, <-chan bool, chan error) {
	inputChannel := make(chan InputRequest)
	outputChannel := make(chan bool)
	terminationChannel := make(chan error, 1)
	for i := 0; i < workers; i++ {
		go work(inputChannel, outputChannel, terminationChannel, graphs)
	}

	return inputChannel, outputChannel, terminationChannel
}
