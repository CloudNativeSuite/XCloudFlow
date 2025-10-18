package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/ssh"
)

// Context provides information for task execution.
type Context struct {
	Host inventory.Host
	Vars map[string]string
	Diff bool
}

// TaskHandler executes a task and returns the result.
type TaskHandler func(ctx Context, task parser.Task) ssh.CommandResult
