package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
	"fmt"
)

func systemdHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Systemd == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing systemd parameters"}
	}
	var before string
	if ctx.Diff {
		before = ssh.RunShellCommand(ctx.Host, fmt.Sprintf("sudo systemctl is-active %s || true", task.Systemd.Name)).Output
	}
	cmd := fmt.Sprintf("sudo systemctl %s %s", task.Systemd.State, task.Systemd.Name)
	if task.Systemd.Enabled {
		cmd = fmt.Sprintf("%s && sudo systemctl enable %s", cmd, task.Systemd.Name)
	}
	res := ssh.RunShellCommand(ctx.Host, cmd)
	if ctx.Diff {
		after := ssh.RunShellCommand(ctx.Host, fmt.Sprintf("sudo systemctl is-active %s || true", task.Systemd.Name)).Output
		res.Output = ssh.Diff(before, after, task.Systemd.Name)
	}
	return res
}

func init() { Register("systemd", systemdHandler) }
