package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func yumHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Yum == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing yum parameters"}
	}
	pkg := task.Yum.Name
	cmd := fmt.Sprintf("sudo yum -y install %s", pkg)
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() { Register("yum", yumHandler) }
