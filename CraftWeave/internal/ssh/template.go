package ssh

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"craftweave/internal/inventory"
)

// RenderTemplate renders the given template file with data and uploads it to the remote host
func RenderTemplate(h inventory.Host, src, dest string, data map[string]string, diff bool) CommandResult {
	content, err := os.ReadFile(src)
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("read template failed: %v", err),
		}
	}

	t, err := template.New("tpl").Parse(string(content))
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("parse template failed: %v", err),
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("execute template failed: %v", err),
		}
	}

	newContent := buf.Bytes()
	var diffText string
	if diff {
		before := RunShellCommand(h, fmt.Sprintf("cat %s 2>/dev/null || true", dest)).Output
		diffText = Diff(before, string(newContent), dest)
	}

	encoded := base64.StdEncoding.EncodeToString(newContent)
	script := fmt.Sprintf("echo \"%s\" | base64 -d > %s", encoded, dest)

	res := RunShellCommand(h, script)
	if diff {
		res.Output = diffText
	}
	return res
}

// UploadFile copies a local file to the remote host at dest path.
func UploadFile(h inventory.Host, src, dest string, diff bool) CommandResult {
	content, err := os.ReadFile(src)
	if err != nil {
		return CommandResult{Host: h.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("read file failed: %v", err)}
	}
	var diffText string
	if diff {
		before := RunShellCommand(h, fmt.Sprintf("cat %s 2>/dev/null || true", dest)).Output
		diffText = Diff(before, string(content), dest)
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	script := fmt.Sprintf("echo \"%s\" | base64 -d > %s", encoded, dest)
	res := RunShellCommand(h, script)
	if diff {
		res.Output = diffText
	}
	return res
}
