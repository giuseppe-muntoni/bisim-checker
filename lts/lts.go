package lts

type Process struct {
	ProcessID string
	Actions   []Action
}

type Action struct {
	ChannelID       string
	Send            bool
	TargetProcessID string
}

type Graph struct {
	ID        string
	Processes map[string]Process
}

func MakeProcess() Process {
	return Process{
		ProcessID: "",
		Actions:   []Action{},
	}
}

func MakeGraph() Graph {
	return Graph{
		Processes: map[string]Process{},
	}
}
