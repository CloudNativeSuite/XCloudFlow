# Xconfig

**Xconfig** 是一个使用 Go 编写的 Ansible-like 工具，支持任务执行、架构编织、图模型导出与插件扩展。

---

## 🧩 特性

- 🛠️ `xconfig remote`：执行单条远程命令（支持 shell 模块）
- 📜 `xconfig playbook`：运行 YAML Playbook
- 🔐 `xconfig vault`：加解密配置(Todo)
- 🧠 `xconfig cmdb`：输出图数据库模型(Todo)
- 🧩 `xconfig plugin`：加载并运行插件（Todo 支持 WASM）

---

## 🚀 快速开始

1. 编译项目 make
2. 执行远程 shell 命令（类似 ansible）
使用 INI 格式的 inventory 文件：

```
[all]
demo           ansible_host=192.168.124.77     ansible_ssh_user=shenlan
cn-hub         ansible_host=1.15.155.245       ansible_ssh_user=ubuntu


[all:vars]
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/id_rsa

3. 执行命令： ./xconfig remote all -i example/inventory -m shell -a 'id'
4. 输出示例：
```
🧶 欢迎使用：Xconfig - 任务与架构编织工具
deepflow-demo | CHANGED | rc=0 >>
uid=1000(shenlan) gid=1000(shenlan) groups=1000(shenlan),10(wheel)

cn-hub | CHANGED | rc=0 >>
uid=1000(ubuntu) gid=1001(ubuntu) groups=1001(ubuntu),27(sudo),
...
```

5. 支持 dry-run 模式：

xconfig remote all -i example/inventory -m shell -a 'id' -C


6. 聚合输出（推荐用于大规模场景）

xconfig remote all -i example/inventory -m shell -a 'id' --aggregate

示例输出：
```
```bash
demo,cn-hub,tky-proxy | CHANGED | rc=0 >>
uid=1000(ubuntu) gid=1000(ubuntu) groups=...

icp-huawei,global-hub | CHANGED | rc=0 >>
uid=0(root) gid=0(root) groups=0(root)


7. 执行本地脚本（上传到远程临时执行）

```bash
chmod +x example/*.sh

./xconfig remote -i example/inventory all -m script -a example/echo.sh
./xconfig remote -i example/inventory all -m script -a example/uname.sh --aggregate
./xconfig remote -i example/inventory all -m script -a example/nproc.sh --aggregate
```bash
```

📌 `--aggregate / -A` 会自动对输出相同的主机进行聚合展示。
```

# ⚙️ 全局参数

参数	描述
--aggregate, -A	聚合输出相同结果的主机，适用于大规模展示

# 📁 项目结构

```
Xconfig/
├── cmd/                  # Cobra 命令定义
│   ├── root.go           # 根命令
│   ├── remote.go        # 类 ansible 子命令
│   ├── playbook.go       # 执行 playbook
│   ├── vault.go          # 加解密相关
│   ├── cmdb.go           # 输出图模型
│   └── plugin.go         # 插件运行
├── core/                 # 核心逻辑模块
│   ├── executor/         # 执行器引擎
│   ├── parser/           # playbook/拓扑解析
│   ├── cmdb/             # 图模型构建与导出
│   └── plugin/           # 插件接口定义与加载
├── internal/             # 内部工具库（如 ssh 执行器、inventory 解析器）
│   └── inventory/
│   └── ssh/
│       ├── result.go       # ➕ 定义 CommandResult
│       ├── formatter.go    # ➕ 实现 AggregatedPrint
│       └── runner.go       # 🔁 改为返回 CommandResult
├── plugins/              # 插件目录（WASM/Go 可选）
├── example/              # 示例 inventory 和脚本
│   ├── inventory         # 测试主机清单
│   ├── echo.sh           # 输出 hostname
│   └── uname.sh          # 输出内核信息
├── banner.txt            # CLI 启动 ASCII 图标
├── go.mod
├── go.sum
├── main.go
└── README.md
```

# 🔮 愿景

Xconfig 旨在成为下一代 DevOps 工具 —— 融合任务调度、架构可视化与智能插件能力，支持轻量化、模块化和智能化的运维体验。

> 辅以 🤖 ChatGPT 之力，愿你我皆成 AIGC 时代的创造者与编织者 🚀
