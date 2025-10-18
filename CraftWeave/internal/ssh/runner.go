package ssh

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"craftweave/internal/inventory"
)

// RunShellCommand 使用 Go 原生 SSH 实现，优先使用私钥，失败时回退密码认证
func RunShellCommand(h inventory.Host, command string) CommandResult {
	var authMethods []ssh.AuthMethod
	var authMethodUsed string

	// 优先尝试私钥认证
	if h.KeyFile != "" {
		if keyBytes, err := os.ReadFile(h.KeyFile); err == nil {
			if signer, err := ssh.ParsePrivateKey(keyBytes); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
				authMethodUsed = "key"
			}
		}
	}

	// 若未成功加载私钥，则尝试密码认证
	if authMethodUsed != "key" && h.Password != "" {
		authMethods = append(authMethods, ssh.Password(h.Password))
		authMethodUsed = "password"
	}

	// 若无有效认证方式，返回失败
	if len(authMethods) == 0 {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     "No valid SSH authentication method found (key or password)",
		}
	}

	config := &ssh.ClientConfig{
		User:            h.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 跳过 host key 校验（生产环境建议自定义）
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", h.Address, h.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("SSH dial error (%s): %v", authMethodUsed, err),
		}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return CommandResult{
			Host:       h.Name,
			ReturnMsg:  "FAILED",
			ReturnCode: 1,
			Output:     fmt.Sprintf("Session error: %v", err),
		}
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)

	result := CommandResult{
		Host:   h.Name,
		Output: string(output),
	}

	if err != nil {
		result.ReturnMsg = "FAILED"
		result.ReturnCode = 1
	} else {
		result.ReturnMsg = "CHANGED"
		result.ReturnCode = 0
	}

	return result
}
