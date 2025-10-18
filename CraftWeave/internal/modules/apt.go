package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func aptHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Apt == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing apt parameters"}
	}
	pkg := task.Apt.Name
	if task.Apt.Deb != "" {
		pkg = task.Apt.Deb
	}
	cmd := fmt.Sprintf("sudo apt-get -y install %s", pkg)
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() { Register("apt", aptHandler) }
