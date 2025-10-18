package ssh

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"craftweave/internal/inventory"
)

// RunRemoteScript uploads a script to a remote host, executes it, and cleans up
func RunRemoteScript(h inventory.Host, scriptPath string) CommandResult {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("read script failed: %v", err),
		}
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	remotePath := fmt.Sprintf("/tmp/craftweave-%d.sh", time.Now().UnixNano())

	script := fmt.Sprintf(`
		echo "%s" | base64 -d > %s && \
		chmod +x %s && \
		(if command -v bash >/dev/null 2>&1; then bash %s; else sh %s; fi); code=$? && \
		rm -f %s && \
		exit $code
	`, encoded, remotePath, remotePath, remotePath, remotePath, remotePath)

	return RunShellCommand(h, script)
}
