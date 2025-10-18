# 🦀 Xconfig Agent

Xconfig Agent (`cw-agent`) 是一个独立运行的本地 Playbook 执行器，支持从 Git 仓库拉取剧本（playbook.yaml），执行本地 shell/script 命令任务。无需 Controller、无需远程 SSH，适用于边缘节点、本地运维任务等场景。

---

## 📁 项目结构``

`
cw-agent/
├── Cargo.toml
├── README.md
├── Makefile
├── example/
│ └── playbook.yaml # 示例本地 playbook
├── cw-agent.service # 可选：systemd 单元文件
└── src/
├── main.rs # 入口，CLI 参数解析 + 调度器入口
├── scheduler.rs # 定时/触发式拉取、执行、保存
├── config.rs # 拉取并解析配置（Git/HTTP、本地文件）
├── executor.rs # 执行任务（shell/script 本地运行）
├── result_store.rs # 存储执行结果（JSON 本地落盘）
└── models.rs # Play / Task 结构体定义
```

## ✅ 功能目标（本地 Playbook 执行器）

- 支持从 Git 仓库拉取 `playbook.yaml`
- 支持任务类型：`shell`、`script`（本地执行）
- 支持 `--oneshot` 或 `daemon` 模式定期执行
- 所有任务限定运行在本机（无 SSH，无 controller）

---

## 🧩 支持命令说明

| 命令格式                      | 功能说明                                               |
|-----------------------------|--------------------------------------------------------|
| `cw-agent oneshot`           | 一次性从 `/etc/cw-agent.conf` 拉取 Git 仓库并执行 Playbook |
| `cw-agent daemon`            | 持续运行，按 interval 定期拉取并执行                   |
| `cw-agent playbook --file x.yaml` | 执行指定本地 Playbook 文件（仅作用于本机）           |
| `cw-agent status`            | 输出最近一次任务执行结果（来自 `/var/lib/cw-agent/`） |
| `cw-agent version`           | 显示版本号信息                                        |

---

## 🧪 示例测试用例

### 1. ✅ 本地运行示例 playbook

```bash
cw-agent playbook --file example/playbook.yaml
内容示例：

yaml
复制
编辑
- name: Local Test
  tasks:
    - name: Print hello
      shell: echo "Hello from Xconfig Agent"

    - name: Show time
      shell: date
2. ✅ 配置 Git 拉取执行
/etc/cw-agent.conf
repo: "https://github.com/your-org/your-repo.git"
interval: 60
playbook:
  - sync/playbook.yaml
示例 playbook.yaml 内容（在 Git 仓库中）
yaml
- name: System Info
  tasks:
    - name: Uptime
      shell: uptime

    - name: Disk usage
      shell: df -h

sudo --preserve-env=HTTPS_PROXY,HOME ./target/release/cw-agent oneshot


# 🛠️ TODO（可选）

- 支持日志落盘与 rotate
- 支持 JSON/YAML 混合格式输入
- 支持 cron 表达式自定义调度
- 支持远程结果上报（WebHook/HTTP）
