# Elastic IAC 设计文档

本文档结合现有项目目录结构，说明如何使用 Go 构建具备弹性、易扩展的基础设施即代码（IAC）框架，并提供与 CMDB、Ansible 以及 AI MCP 协议集成的思路。

## 1. 项目目录回顾

项目遵循 Go Module 和 `cobra` 命令行结构：

```
xcloud-cli/
├── cmd/               # CLI 命令入口
├── internal/          # 内部实现，避免外部引用
│   └── pulumi/        # Pulumi 基础设施代码
├── docs/              # 设计文档
├── main.go            # 程序入口
└── go.mod
```

该结构方便扩展新的功能模块，例如对接 CMDB、Ansible 等。

## 2. 设计目标

1. **弹性与扩展**：通过模块化 Go 代码便于后续接入更多云平台或新的运维流程。
2. **CMDB 集成**：在部署前后与 CMDB 同步资源信息，实现资产管理闭环。
3. **Ansible 支持**：在基础设施部署完成后调用 Ansible 对实例进行配置。
4. **AI MCP 协议**：预留 AI 管控接口，便于未来通过 MCP 协议实现智能调度或自动化运维。

## 3. 关键模块划分

- `internal/pulumi`：封装 Pulumi 资源定义与部署逻辑。
- `internal/cmdb`：提供与 CMDB 系统交互的接口，如资源登记、状态更新等。
- `internal/ansible`：调用 Ansible Playbook 的封装，保持与现有 `cmd/ansible.go` 兼容。
- `internal/mcp`：示例化 AI MCP 协议交互的基础代码。

## 4. 基础代码结构示例

以下示例展示了新的内部包结构以及主要函数定义。

```go
// internal/cmdb/cmdb.go
package cmdb

// Resource 表示需要登记到 CMDB 的基础结构。
type Resource struct {
    ID   string
    Type string
    Tags map[string]string
}

// Register 用于在资源创建后向 CMDB 注册信息。
func Register(r Resource) error {
    // TODO: 调用实际的 CMDB 接口
    return nil
}
```

```go
// internal/ansible/run.go
package ansible

import "os/exec"

// RunPlaybook 调用 ansible-playbook 执行配置。
func RunPlaybook(inventory, playbook string) error {
    cmd := exec.Command("ansible-playbook", "-i", inventory, playbook)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
```

```go
// internal/mcp/mcp.go
package mcp

// Send 通过 MCP 协议下发指令或接收状态。
func Send(payload []byte) error {
    // TODO: 实现与 AI MCP 系统的通信
    return nil
}
```

在 `cmd/` 下可按需增加新命令，例如 `cmd/cmdb.go` 或 `cmd/mcp.go`，并在 `root.go` 中注册。

## 5. 部署流程示例

1. 运行 `xcloud-cli up` 部署基础设施。
2. 部署成功后调用 `cmdb.Register` 将资源信息同步到 CMDB。
3. 使用 `ansible.RunPlaybook` 对新建实例执行后续配置。
4. 若需要 AI 调度，可通过 `mcp.Send` 将状态汇报给 MCP 系统。

## 6. 未来规划

- 丰富 `internal/mcp` 实现，使其支持消息队列或 WebSocket 等实时协议。
- 增加更多云资源模块，与 `internal/pulumi` 协作输出统一的资源结构，便于 CMDB 登记。
- 在 CI/CD 流程中引入安全扫描与合规检查，保证整体流程稳定可靠。

