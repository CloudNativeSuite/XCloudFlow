package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func statHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Stat == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing stat parameters"}
	}
	cmd := fmt.Sprintf("test -e %s && echo exists || echo missing", task.Stat.Path)
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() { Register("stat", statHandler) }
