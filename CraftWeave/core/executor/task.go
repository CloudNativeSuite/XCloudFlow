package executor

import (
	"fmt"

	"craftweave/core/parser"
	"craftweave/internal/inventory"
	"craftweave/internal/modules"
	"craftweave/internal/ssh"
)

// ExecuteTask dispatches the task to the appropriate module handler.
func ExecuteTask(task parser.Task, host inventory.Host, vars map[string]string, diff bool) ssh.CommandResult {
	ctx := modules.Context{Host: host, Vars: vars, Diff: diff}

	var res ssh.CommandResult
	if h, ok := modules.GetHandler(task.Type()); ok {
		res = h(ctx, task)
	} else {
		switch {
		case task.Shell != "":
			res = ssh.RunShellCommand(host, task.Shell)
		case task.Script != "":
			res = ssh.RunRemoteScript(host, task.Script)
		case task.Template != nil:
			res = ssh.RenderTemplate(host, task.Template.Src, task.Template.Dest, vars, diff)
		default:
			res = ssh.CommandResult{Host: host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("Unsupported task type in '%s'", task.Name)}
		}
	}

	if task.Register != "" {
		vars[task.Register] = res.Output
	}
	return res
}
