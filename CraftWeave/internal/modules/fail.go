package modules

import (
	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func failHandler(ctx Context, task parser.Task) ssh.CommandResult {
	msg := ""
	if task.Fail != nil {
		msg = task.Fail.Msg
	}
	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: msg}
}

func init() { Register("fail", failHandler) }
