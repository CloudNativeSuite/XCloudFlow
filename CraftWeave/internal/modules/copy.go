package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func copyHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.Copy == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing copy parameters"}
	}
	return ssh.UploadFile(ctx.Host, task.Copy.Src, task.Copy.Dest, ctx.Diff)
}

func init() { Register("copy", copyHandler) }
