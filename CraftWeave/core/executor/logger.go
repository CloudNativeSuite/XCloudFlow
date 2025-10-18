package executor

import "craftweave/internal/ssh"

// LogCollector collects command execution results.
type LogCollector interface {
	Collect(ssh.CommandResult)
}
