package modules

import (
	"bytes"
	"text/template"

	"craftweave/core/parser"
	"craftweave/internal/ssh"
)

func shellHandler(ctx Context, task parser.Task) ssh.CommandResult {
	cmd := task.Shell
	if len(ctx.Vars) > 0 {
		if tmpl, err := template.New("shell").Parse(task.Shell); err == nil {
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, ctx.Vars); err == nil {
				cmd = buf.String()
			}
		}
	}
	return ssh.RunShellCommand(ctx.Host, cmd)
}

func init() {
	Register("shell", shellHandler)
}
