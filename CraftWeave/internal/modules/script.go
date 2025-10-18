package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func scriptHandler(ctx Context, task parser.Task) ssh.CommandResult {
	return ssh.RunRemoteScript(ctx.Host, task.Script)
}

func init() {
	Register("script", scriptHandler)
}
