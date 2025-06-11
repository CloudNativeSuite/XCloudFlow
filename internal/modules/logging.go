package modules

import "sync"

// LogCollector receives output from executed tasks.
type LogCollector interface {
	Collect(taskType, output string)
}

var (
	collector LogCollector = &stdoutCollector{}
)

// SetCollector sets a global log collector.
func SetCollector(c LogCollector) {
	if c != nil {
		collector = c
	}
}

func dispatchOutput(taskType, output string) {
	if collector != nil && output != "" {
		collector.Collect(taskType, output)
	}
}

type stdoutCollector struct {
	mu sync.Mutex
}

func (s *stdoutCollector) Collect(taskType, output string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	println("["+taskType+"]", output)
}
