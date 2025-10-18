# Git Repo Provider 设计

本文档规划一个统一的 Git 仓库 Provider，用于在 Pulumi 等 IaC 场景下管理各类 Git 托管服务。目标是同时支持 GitHub、GitLab 以及 GitHub-like 的本地部署（如 Gitea、Forgejo），并提供以下能力：

## 1. 核心目标

1. **统一接口**：抽象常用仓库能力，屏蔽不同平台间的 API 差异。
2. **可扩展性**：通过插件或适配层扩展新的 Git 托管实现。
3. **声明式配置**：结合 Pulumi/Terraform 等工具，使用代码描述仓库及其策略。

## 2. 功能范围

### 2.1 Workflow 配置
- 支持导入/同步 CI 配置（GitHub Actions、GitLab CI、Gitea Actions 等）。
- 允许启用/禁用某些 workflow，并管理默认分支策略。

### 2.2 分支保护与规则集
- **通用分支保护**：只读/禁止强制推送、最少审查人、签名提交等。
- **规则集（Rulesets）**：针对符合模式的分支应用更精细的策略；示例：`release/*` 必须通过 PR、线性历史、状态检查等。
- 支持各平台特有功能：
  - GitHub RepositoryRuleset
  - GitLab Branch Protection
  - 本地实例的相应 API。

## 3. 架构草图

```
+-----------------------+
| Git Repo Provider SDK |
+-----------------------+
          |
          +-- Adapter: GitHub
          |
          +-- Adapter: GitLab
          |
          +-- Adapter: Self-hosted (Gitea 等)
```

- 核心 SDK 提供统一的 `Repository`、`BranchRule`、`Workflow` 等资源定义。
- 每个 Adapter 负责将统一定义转换为目标平台的 API 调用。
- 通过接口或 `go build` 标签选择具体实现。

## 4. Pulumi 示例

与本仓库下的 `infra/main.go` 类似，可在不同平台传入相应 Provider：

```go
// 以 GitHub 为例
_, err := github.NewRepositoryRuleset(ctx, "protect-release-pattern", &github.RepositoryRulesetArgs{ ... })
```

针对 GitLab 或本地实例，可替换为对应的 Pulumi Provider 或自定义组件。

## 5. 后续工作

1. 设计统一的资源描述（YAML/Go Struct）。
2. 实现 GitHub/GitLab Adapter，并验证基础功能。
3. 接入自托管服务（Gitea/Forgejo），确保 API 兼容。
4. 提供示例与文档，方便用户使用。

