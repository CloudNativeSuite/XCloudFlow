# ☁️ XCloudFlow Monorepo

XCloudFlow 将多云基础设施、配置编排与边缘执行整合在一个仓库中。仓库内包含三个相互独立、又能协同工作的 CLI：

| 目录            | 语言 | 说明 |
|-----------------|------|------|
| `xcloud-cli/`   | Go   | 面向 Terraform/Pulumi 场景的控制平面 CLI，统一管理多云部署生命周期。|
| `xconfig/`      | Go   | 类 Ansible 的任务/剧本执行器，提供 `remote`/`playbook` 等命令。|
| `xconfig-agent/`| Rust | 轻量级边缘 Agent，周期性拉取剧本并在本地执行，支撑无人值守环境。|

---

## 🚀 快速开始

### 1. xcloud-cli（控制平面 CLI）
```bash
cd xcloud-cli
make build        # or `go run main.go --env sit up`
```
常用子命令：`up`、`down`、`export`、`import`、`ansible`。详情见 `xcloud-cli/Makefile` 或 `xcloud-cli/cmd/*.go`。

### 2. xconfig（任务编排 CLI）
```bash
cd xconfig
make build
./xconfig remote all -i example/inventory -m shell -a 'id'
```
- `xconfig remote`：远程命令执行（shell/command/copy/service 等模块）。
- `xconfig playbook`：运行 YAML Playbook，支持 `template`、`setup`、`apt/yum` 等模块。
- 更多示例参见 `xconfig/example/` 与 `xconfig/README.md`。

### 3. xconfig-agent（边缘执行 Agent）
```bash
cd xconfig-agent
cargo build --release
./target/release/xconfig-agent oneshot
```
默认配置从 `/etc/xconfig-agent.conf` 拉取 Git 仓库、读取 Playbook，并将执行结果落盘到 `/var/lib/xconfig-agent/`。

---

## 🧰 仓库级 Makefile

根目录提供一份聚合 `Makefile`，可快速调用各子项目命令：
```bash
make help
make xcloud-build
make xconfig-playbook
make xconfig-agent-run
```

---

## 📚 设计文档

详见 `docs/` 目录，涵盖：
- `XCloudFlowDesign.md`：整体平台架构
- `ModuleExecutionDesign.md`：模块化执行框架设计
- `ElasticIACDesign.md`：Go + Pulumi 弹性 IAC 架构
- `craftweave-playbook-spec.md`：Xconfig Playbook DSL

---

## 🤝 贡献

1. Fork 并创建功能分支。
2. 在对应子目录内运行 `make test`/`cargo test`（如适用）。
3. 提交 PR 并附上测试记录。

欢迎提出 Issue 或 PR，一起打造云管 + 配置 + 边缘执行的一体化工作流。☁️🧵🦀
