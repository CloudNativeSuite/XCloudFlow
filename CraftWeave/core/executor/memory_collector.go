package executor

import (
	"craftweave/internal/ssh"
	"sync"
)

// MemoryCollector stores command results in memory.
type MemoryCollector struct {
	mu      sync.Mutex
	Results []ssh.CommandResult
}

// Collect appends a result to the in-memory slice.
func (m *MemoryCollector) Collect(res ssh.CommandResult) {
	m.mu.Lock()
	m.Results = append(m.Results, res)
	m.mu.Unlock()
}
