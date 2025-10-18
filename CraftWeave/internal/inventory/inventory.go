package inventory

import (
	"bufio"
	"os"
	"strings"
)

type Host struct {
	Name     string
	Address  string
	User     string
	KeyFile  string
	Port     string
	Password string // ✅ 新增：支持密码登录
}


func Parse(path, group string) ([]Host, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []Host
	scanner := bufio.NewScanner(file)
	currentGroup := ""
	defaultUser := "ubuntu"
	defaultPort := "22"
	keyFile := os.Getenv("HOME") + "/.ssh/id_rsa"

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentGroup = strings.Trim(line, "[]")
			continue
		}
		if currentGroup == group {
			parts := strings.Fields(line)
			h := Host{Name: parts[0], Address: parts[0], User: defaultUser, Port: defaultPort, KeyFile: keyFile}
			for _, p := range parts[1:] {
				if strings.HasPrefix(p, "ansible_host=") {
					h.Address = strings.TrimPrefix(p, "ansible_host=")
				}
				if strings.HasPrefix(p, "ansible_ssh_user=") {
					h.User = strings.TrimPrefix(p, "ansible_ssh_user=")
				}
				if strings.HasPrefix(p, "ansible_port=") {
					h.Port = strings.TrimPrefix(p, "ansible_port=")
				}
				if strings.HasPrefix(p, "ansible_ssh_private_key_file=") {
					h.KeyFile = strings.TrimPrefix(p, "ansible_ssh_private_key_file=")
				}
			}
			hosts = append(hosts, h)
		}
	}
	return hosts, nil
}

