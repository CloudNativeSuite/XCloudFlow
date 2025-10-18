# Git Repo Provider 子任务拆分

基于 [git-repo-provider.md](git-repo-provider.md) 的设计，以下列出实现统一 Git 仓库 Provider 所需的主要子任务，按功能与平台分类，便于分工与跟踪进度。

## 1. 统一资源定义
- 设计 `Repository`、`Workflow`、`BranchRule` 等通用结构。
- 提供 YAML/Go Struct 表达，支持 Pulumi/Terraform 映射。
- 规划跨平台字段与平台特有字段的扩展方式。

## 2. Adapter 实现
### 2.1 GitHub
- 封装仓库创建、更新与合并策略 API。
- 同步 GitHub Actions workflow（启用/禁用、默认分支）。
- 实现 RepositoryRuleset：PR 合并、线性历史、状态检查等。

### 2.2 GitLab
- 映射仓库与分支保护接口。
- 管理 `.gitlab-ci.yml` workflow 与默认分支。
- 适配 GitLab Branch Protection 的审查、强推等策略。

### 2.3 自托管 (Gitea/Forgejo)
- 接入相应 REST API，验证版本兼容性。
- 支持 Actions/runner 配置同步。
- 提供分支保护与受限推送规则。

## 3. 工作流与规则集管理
- 定义启用/禁用 workflow 的通用流程。
- 统一分支保护与规则集配置接口，允许按模式（如 `release/*`）应用。
- 支持自定义状态检查与审批人数等高级策略。

## 4. 示例与文档
- 在 `infra/` 提供 GitHub、GitLab、自托管三类示例。
- 编写使用说明与 README 更新。
- 记录平台差异与限制，提供故障排查建议。

## 5. 测试与发布
- 为各 Adapter 编写单元测试与集成测试。
- 设置 CI 以运行测试并验证 Pulumi 预览。
- 发布初始版本并收集反馈。

