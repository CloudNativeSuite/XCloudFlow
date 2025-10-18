package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func templateHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Template == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "template missing"}
	}
	return ssh.RenderTemplate(ctx.Host, task.Template.Src, task.Template.Dest, ctx.Vars, ctx.Diff)
}

func init() {
	Register("template", templateHandler)
}
