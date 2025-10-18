package executor

import "xconfig/internal/ssh"

// LogCollector collects command execution results.
type LogCollector interface {
	Collect(ssh.CommandResult)
}
