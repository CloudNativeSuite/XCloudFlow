# 🧶 Xconfig

**Xconfig** 是一个融合任务编排、架构建模和轻量执行能力的DevOps 工具，灵感源于 Ansible，但更灵活、模块化，支持图模型导出与插件扩展，并具备独立运行的 Rust 轻量 Agent。

---

## 🧩 特性概览

### ✅ 控制端（Go 实现）
- 🛠️ `xconfig remote`：类 Ansible 远程命令执行（内置 shell、command、copy、service 等模块）
- 📜 `xconfig playbook`：YAML 多步骤任务编排（支持 template、setup、apt/yum 等模块）
- 🔐 `xconfig vault`：加解密配置 (TODO)
- 🧠 `xconfig cmdb`：导出拓扑图数据库 (TODO)
- 🔌 `xconfig plugin`：支持插件执行，预留 WASM 接口 (TODO)

### ✅ Agent 端（Rust 实现）
- 🧩 定时拉取配置并执行命令/copy/service 等操作
- 📦 独立运行、轻量部署，可在任意主机常驻执行

---

## 🚀 快速开始（控制端 CLI）

1. 编译 make
2. 示例 inventory 文件（INI 格式）

```
[all]
demo           ansible_host=192.168.124.77     ansible_ssh_user=shenlan role=demo
cn-hub         ansible_host=1.15.155.245       ansible_ssh_user=ubuntu
global-hub     ansible_host=2.15.135.215       ansible_ssh_user=centos

[all:vars]
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/id_rsa
env='prod'
```

3. 远程执行命令（Shell 模块）

xconfig remote all -i example/inventory -m shell -a 'id'

4. 远程执行命令（Command 模块）

xconfig remote all -i example/inventory -m command -a '/usr/bin/id'

5. 输出聚合展示（推荐用于大规模场景

xconfig remote all -i example/inventory -m shell -a 'id' --aggregate

6. 上传并执行脚本

xconfig remote all -i example/inventory -m script -a example/uname.sh

7. Dry-run 模式预览执行效果

xconfig remote all -i example/inventory -m shell -a 'id' -C

8. 显示文件或模板差异（可与 -A、-C 组合）

xconfig remote all -i example/inventory -m template -a "example.j2:/tmp/out.yml" -D -C -A

9. 指定单个主机运行命令

xconfig remote cn-hub -i example/inventory -m shell -a 'uptime'

10. 聚合输出执行脚本结果

xconfig remote all -i example/inventory -m script -a example/nproc.sh --aggregate

11. 运行 Playbook 文件

xconfig playbook -i example/inventory example/run_example -i example/inventory

可选：执行更复杂的示例

- xconfig playbook -i example/inventory example/playbooks/system-check.yaml 
- xconfig playbook -i example/inventory example/playbooks/set-password.yml -e password=YOURPASS
- xconfig playbook -i example/inventory example/deploy_deepflow_agent

# 📦 Agent 支持命令说明

| 命令格式                      | 功能说明                                               |
|-----------------------------|--------------------------------------------------------|
| `cw-agent oneshot`           | 一次性从 `/etc/cw-agent.conf` 拉取 Git 仓库并执行 Playbook |
| `cw-agent daemon`            | 持续运行，按 interval 定期拉取并执行                   |
| `cw-agent playbook --file x.yaml` | 执行指定本地 Playbook 文件（仅作用于本机）           |
| `cw-agent status`            | 输出最近一次任务执行结果（来自 `/var/lib/cw-agent/`） |
| `cw-agent version`           | 显示版本号信息

# ⚙️ 全局参数

| 参数              | 描述                                               |
|-------------------|----------------------------------------------------|
| `--aggregate`, `-A` | 聚合输出相同结果的主机（大规模场景推荐）         |
| `--check`, `-C`     | Dry-run 模式，不实际执行命令（TODO）              |
| `--diff`, `-D`      | 当修改文件和模板时显示差异                         |
| `--extra-vars`, `-e` | 运行时变量，覆盖 Playbook 中的 `vars`             |

# 控制端（Go 实现）

📁 项目结构
```
Xconfig/
├── cmd/                  # CLI 命令定义（Cobra）
├── core/                 # 核心执行器、解析器、拓扑建模
├── internal/             # SSH 库、Inventory 处理
├── plugins/              # 插件接口定义与运行（TODO）
├── example/              # 示例 inventory + 脚本
├── banner.txt            # CLI 欢迎图标
├── XconfigAgent/      # Rust 版 cw-agent
│   └── src/              # agent 源码目录
└── main.go
```

#🧠 Xconfig Agent（Rust 实现）

📦 结构目录
```
Xconfig-agent/
├── Cargo.toml
├── cw-agent.service              # systemd 单元文件（可选）
└── src/
    ├── main.rs                  # CLI 入口，启动/守护/状态
    ├── scheduler.rs             # 定时拉取与调度执行
    ├── config.rs                # 解析配置（JSON、Git、HTTP）
    ├── executor.rs              # 支持 command/copy/service 执行
    ├── result_store.rs          # 本地 JSON/DB 结果保存
    └── models.rs                # 配置/结果结构体定义
```

# 🔮 愿景

Xconfig 旨在成为一个轻量级 DevOps 工具，融合任务调度、配置编排、架构建模、图数据库与 AI 辅助的智能插件系统，构建“人-机-架构”高效协作闭环。

> 借助 🤖 ChatGPT 之力，愿你我皆成 AIGC 时代的创造者与织梦者 🚀
