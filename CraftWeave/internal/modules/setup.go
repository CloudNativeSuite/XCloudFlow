package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func setupHandler(ctx Context, task parser.Task) ssh.CommandResult {
	res := ssh.RunShellCommand(ctx.Host, "uname -a")
	ctx.Vars["ansible_facts"] = res.Output
	return res
}

func init() {
	Register("setup", setupHandler)
	Register("gather_facts", setupHandler)
}
