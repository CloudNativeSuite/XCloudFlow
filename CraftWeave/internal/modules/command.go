package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func commandHandler(ctx Context, task parser.Task) ssh.CommandResult {
	return ssh.RunShellCommand(ctx.Host, task.Command)
}

func init() { Register("command", commandHandler) }
