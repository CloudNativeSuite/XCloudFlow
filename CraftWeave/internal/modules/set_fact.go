package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func setFactHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.SetFact != nil {
		for k, v := range task.SetFact {
			ctx.Vars[k] = v
		}
	}
	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "OK", ReturnCode: 0, Output: ""}
}

func init() { Register("set_fact", setFactHandler) }
