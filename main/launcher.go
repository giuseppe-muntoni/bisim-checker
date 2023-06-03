package main

import (
	"bisim-checker/checker"
	"bisim-checker/io"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: " + os.Args[0] + " first_graph.dot second_graph.dot number_of_workers")
		return
	}

	firstGraphPath := os.Args[1]
	secondGraphPath := os.Args[2]

	parDegree, err := strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if bisimulation, err := io.CreateBisimulationFromFiles(firstGraphPath, secondGraphPath).Get(); err != nil {
		fmt.Println(err.Error())
	} else {
		if err := checker.Emit(bisimulation, parDegree); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Success!")
		}
	}
}
