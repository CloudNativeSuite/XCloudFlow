# XCloudFlow 状态服务存储设计

本文档按照“Runner(TF/Pulumi) ──HTTP──► 状态服务 ──► PostgreSQL”的体系结构，说明在 XCloudFlow 中如何基于 PostgreSQL 构建状态一致、易扩展的存储方案，并支持向 CMDB、工单、监控等系统进行 SQL 级联。

## 1. 总体架构

```
Runner(TF/Pulumi) ──HTTP──► 状态服务
                            │
                            ▼
                        PostgreSQL
                   ┌────────┴────────┐
                   │  状态(版本/锁)   │  ← 事务一致
                   │  资源事件/快照   │
                   └────────┬────────┘
                            ▼
             SQL 级联：CMDB / 工单 / 监控
```

- **Runner** 负责执行 Terraform、Pulumi 等工具，并通过 HTTP API 将状态变更、执行结果推送给状态服务。
- **状态服务** 是 XCloudFlow 的核心后端，提供状态管理、锁控制、事件追踪等能力。
- **PostgreSQL** 承载所有结构化数据，通过事务与行级锁保证 Runner 之间的互斥和一致性。
- **SQL 级联** 通过数据库触发器、视图或逻辑订阅，向 CMDB、工单、监控等外部系统同步数据。

## 2. 数据库命名空间

所有表位于 `state` schema 下，按不同职责划分：

- `state.runs`：记录 Runner 发起的执行任务。
- `state.states`：存储基础设施状态（版本、锁、序列号）。
- `state.resources`：以最新快照形式存储资源属性。
- `state.resource_events`：记录资源级事件流，方便回放与审计。
- `state.notifications`：向外部系统推送的消息队列表，通过 SQL 级联生成。

## 3. 核心表设计

### 3.1 state.runs

| 字段             | 类型          | 说明                       |
|------------------|---------------|----------------------------|
| `id`             | UUID          | 运行唯一标识               |
| `runner_type`    | TEXT          | `terraform`、`pulumi` 等   |
| `module`         | TEXT          | 模块/项目名称              |
| `status`         | TEXT          | `pending`、`applying` 等   |
| `started_at`     | TIMESTAMPTZ   | 开始时间                   |
| `finished_at`    | TIMESTAMPTZ   | 结束时间                   |
| `log_url`        | TEXT          | 日志或 Artifact 链接       |

该表用于查询执行历史，与锁表协同避免同一模块并发运行。

### 3.2 state.states

| 字段             | 类型        | 说明                                         |
|------------------|-------------|----------------------------------------------|
| `id`             | UUID        | 状态记录 ID                                  |
| `module`         | TEXT        | 对应模块/环境                               |
| `version`        | BIGINT      | 版本号，自增                                 |
| `lock_owner`     | TEXT        | 持有锁的 Runner 标识                         |
| `lock_acquired`  | TIMESTAMPTZ | 锁获取时间                                   |
| `state_blob`     | JSONB       | Terraform/Pulumi 状态内容                    |
| `checksum`       | TEXT        | 状态文件哈希，用于并发控制                  |
| `updated_at`     | TIMESTAMPTZ | 最近更新时间                                 |

- 通过 `module` + `version` 唯一索引保证版本线性增长。
- 获取锁时使用 `SELECT ... FOR UPDATE`，并在事务结束后释放。

### 3.3 state.resources

| 字段            | 类型        | 说明                               |
|-----------------|-------------|------------------------------------|
| `id`            | UUID        | 资源唯一标识                       |
| `module`        | TEXT        | 所属模块/环境                      |
| `urn`           | TEXT        | Terraform Address / Pulumi URN     |
| `type`          | TEXT        | 资源类型，如 `aws_instance`        |
| `attributes`    | JSONB       | 资源最新属性快照                   |
| `state_version` | BIGINT      | 对应的状态版本                     |
| `updated_at`    | TIMESTAMPTZ | 最新更新时间                       |

该表通过覆盖索引 (`module`, `urn`) 支持快速查询资源当前状态。

### 3.4 state.resource_events

| 字段            | 类型        | 说明                               |
|-----------------|-------------|------------------------------------|
| `id`            | BIGSERIAL   | 自增主键                           |
| `resource_id`   | UUID        | 对应 `resources.id`                |
| `run_id`        | UUID        | 对应 `runs.id`                     |
| `action`        | TEXT        | `create`、`update`、`delete` 等    |
| `diff`          | JSONB       | 属性差异                           |
| `occurred_at`   | TIMESTAMPTZ | 事件发生时间                       |

- 事件表遵循追加模式，方便回放和审计。
- 可以基于 `run_id` 或 `resource_id` 建立分区提升查询性能。

### 3.5 state.notifications

| 字段            | 类型        | 说明                               |
|-----------------|-------------|------------------------------------|
| `id`            | BIGSERIAL   | 自增主键                           |
| `resource_id`   | UUID        | 对应 `resources.id`                |
| `module`        | TEXT        | 模块名称，便于下游过滤             |
| `event_type`    | TEXT        | `cmdb_sync`、`ticket` 等            |
| `payload`       | JSONB       | 下游所需的序列化数据               |
| `created_at`    | TIMESTAMPTZ | 生成时间                           |
| `processed_at`  | TIMESTAMPTZ | 下游处理完成时间                   |

该表可被外部消费程序轮询，也可以通过逻辑复制流推送到消息队列。

## 4. 事务与锁控制

1. Runner 在执行前通过 `BEGIN` + `SELECT FOR UPDATE` 方式占用 `state.states` 表的锁，写入 `lock_owner` 与 `lock_acquired`。
2. 如果锁已存在且未超时，状态服务返回 423 Locked，Runner 按策略重试。
3. 执行成功后，Runner 将新的状态 JSONB 与版本号写入 `state.states`，并在同一事务中更新 `state.resources` 与 `state.resource_events`。
4. 事务提交后通过数据库触发器将变化写入 `state.notifications`，供外部系统消费。

## 5. SQL 级联与外部系统集成

- **CMDB**：创建触发器 `AFTER INSERT OR UPDATE ON state.resources`，将资源属性映射到 CMDB 对应的资产表。
- **工单系统**：通过 `state.notifications` 表生成包含差异信息的 JSON，供自动化工单流程拉取。
- **监控告警**：基于 `resource_events` 构建物化视图，将关键字段映射到 Prometheus Exporter 或告警规则。

外部系统可使用以下方式接入：

1. 直接通过 FDW（Foreign Data Wrapper）在 CMDB 数据库中引用 `state` schema。
2. 使用 Debezium/Logical Replication 订阅 `state.notifications`，推送到 Kafka/RabbitMQ。
3. 定期执行 SQL 视图，将资源快照同步到报表或监控数据仓库。

## 6. 备份与审计

- 使用 PostgreSQL 原生的 `pg_dump` 与 PITR（Point-In-Time Recovery）保证状态可恢复。
- `resource_events` 与 `runs` 表提供完整审计轨迹，可在合规场景中回溯变更。
- 开启 `pgcrypto` 或 `row level security` 保护敏感字段，如凭据、密钥等。

## 7. 性能与扩展

- 根据模块或项目对 `state.resources`、`state.resource_events` 进行分区，降低单表规模。
- 为 JSONB 字段建立 GIN 索引，提升属性查询效率。
- 通过连接池（如 PgBouncer）提升状态服务并发处理能力。
- 若需要跨区域高可用，可利用 PostgreSQL 流复制或云厂商托管服务。

## 8. 迁移与演进

- 使用 `golang-migrate` 或 `atlas` 管理 schema 版本，确保 CICD 中自动执行迁移。
- 为不同 Runner 类型预留扩展字段，如 Helm/Kustomize 等。
- 可以引入 `state.snapshots` 表保存每个版本的压缩归档，以便快速回滚或在调试时下载。

以上设计确保 XCloudFlow 在多 Runner、多环境场景下具备可靠的状态存储与外部系统联动能力。
