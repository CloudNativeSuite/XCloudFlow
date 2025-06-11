package ansible

import (
	"os"
	"os/exec"
)

// RunPlaybook 调用 ansible-playbook 执行配置。
func RunPlaybook(inventory, playbook string) error {
	cmd := exec.Command("ansible-playbook", "-i", inventory, playbook)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
