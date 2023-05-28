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
	Processes map[string]Process
}
