package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func serviceHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Service == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing service parameters"}
	}
	var before string
	if ctx.Diff {
		before = ssh.RunShellCommand(ctx.Host, fmt.Sprintf("sudo service %s status || true", task.Service.Name)).Output
	}
	cmd := fmt.Sprintf("sudo service %s %s", task.Service.Name, task.Service.State)
	if task.Service.Enabled {
		cmd = fmt.Sprintf("%s && sudo systemctl enable %s", cmd, task.Service.Name)
	}
	res := ssh.RunShellCommand(ctx.Host, cmd)
	if ctx.Diff {
		after := ssh.RunShellCommand(ctx.Host, fmt.Sprintf("sudo service %s status || true", task.Service.Name)).Output
		res.Output = ssh.Diff(before, after, task.Service.Name)
	}
	return res
}

func init() { Register("service", serviceHandler) }
