package modules

import (
	"context"
	"fmt"
	"sync"
)

var (
	registryMu sync.RWMutex
	registry   = make(map[string]Handler)
)

// Register registers a handler with a task type key.
func Register(taskType string, h Handler) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[taskType] = h
}

// GetHandler returns the handler registered for the given task type.
func GetHandler(taskType string) (Handler, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	h, ok := registry[taskType]
	return h, ok
}

// ExecuteTask finds the handler by type and executes it.
func ExecuteTask(ctx context.Context, task Task) error {
	h, ok := GetHandler(task.Type())
	if !ok {
		return fmt.Errorf("no handler registered for %s", task.Type())
	}
	out, err := h.Run(ctx, task)
	dispatchOutput(task.Type(), out)
	return err
}
