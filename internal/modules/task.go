package modules

import "context"

// Task defines an executable task with a type identifier.
type Task interface {
	Type() string
}

// Handler executes a task and returns optional output.
type Handler interface {
	Run(ctx context.Context, task Task) (string, error)
}
