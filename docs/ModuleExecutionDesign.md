# 模块化执行框架设计

本文档描述在 `xcloud` 项目中引入可插拔模块、并发控制及日志/CMDB 集成的方案，满足以下特性：

1. **模块注册机制**：任务模块实现统一接口，通过 `modules.Register` 完成注册。
2. **解耦执行**：执行器根据任务类型查找处理器 `modules.ExecuteTask`，避免调用方与具体实现耦合。
3. **并发控制**：`modules.Pool` 控制同时运行的 goroutine 数量，防止资源耗尽。
4. **输出分发**：定义 `modules.LogCollector` 接口，可接入不同日志收集后端。
5. **资源状态写入 CMDB**：通过 `modules.CMDB` 接口支持将 IAC 资源信息同步到图数据库或导出文件。

## 目录结构

```
internal/modules/      # 通用框架代码
    task.go            # Task/Handler 接口定义
    registry.go        # 模块注册与执行
    pool.go            # 并发池实现
    logging.go         # 日志收集器接口及默认实现
    cmdb.go            # CMDB 后端接口
internal/pulumi/
    module.go          # Pulumi 部署任务模块示例
    infra.go           # 具体的 Pulumi 资源定义
```

## 时序示例

1. 具体模块在 `init()` 中调用 `modules.Register` 完成注册。
2. 外部调用创建实现 `modules.Task` 的任务实例。
3. 通过 `modules.ExecuteTask` 执行，该函数根据任务类型找到对应 `Handler`。
4. 执行结果通过 `LogCollector` 统一输出，同时可调用 `modules.RecordResource` 将资源状态写入 CMDB。
5. 多个任务可由 `modules.Pool` 提交并发运行，`Wait()` 在需要时阻塞等待全部完成。

该设计使新增任务模块仅需实现 `Handler` 接口并在 `init()` 注册，即可被 CLI 或其他调度逻辑调用，实现高度可扩展的自动化框架。

## Pulumi 运行结果结构

`internal/pulumi/infra.go` 中的 `DeployInfrastructure` 会返回结构化的 `DeploymentResult`，包含：

- `activeEnv`：当前执行的环境（例如 dev/sit/prod）。
- `configSources`：参与合并的配置目录列表，方便追踪来源。
- `targets`：矩阵展开后的每个 stack 结果，元素为 `cloud`/`region`/`stack`/`status` 等字段。
  - 当 `status=applied` 时，附带 `planPreview`、`applied`、`costEstimate`、`outputs`、`artifacts`、`previewUrl`、`updateUrl`。
  - 当云厂商暂未支持时，返回 `status=skipped` 并给出原因。

`deployHandler.Run` 会将 `DeploymentResult` 序列化为 JSON，透过框架的 `LogCollector` 返回给调用方，从而满足“资源规划预览 / 成本预估 / 交付物与输出”统一回传的需求，并兼容多云矩阵场景。

