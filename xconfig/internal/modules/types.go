package modules

import (
	"xconfig/core/parser"
	"xconfig/internal/inventory"
	"xconfig/internal/ssh"
)

// Context provides information for task execution.
type Context struct {
	Host inventory.Host
	Vars map[string]interface{}
	Diff bool
}

// TaskHandler executes a task and returns the result.
type TaskHandler func(ctx Context, task parser.Task) ssh.CommandResult
