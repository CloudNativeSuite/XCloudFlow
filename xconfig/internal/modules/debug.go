package modules

import (
	"xconfig/core/parser"
	"xconfig/internal/ssh"
)

func debugHandler(ctx Context, task parser.Task) ssh.CommandResult {
	msg := ""
	if task.Debug != nil {
		msg = task.Debug.Msg
	}
	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "OK", ReturnCode: 0, Output: msg}
}

func init() { Register("debug", debugHandler) }
