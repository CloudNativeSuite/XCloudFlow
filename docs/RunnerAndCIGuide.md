# Runner 与 CI 集成指南

本文档说明如何在本地执行 Pulumi Runner，并通过 API 触发 GitHub Actions/GitLab CI 流水线。

## 1. 本地 Runner 运行

1. 选择环境配置：示例目录 `example/config` 已拆分为 `base/` + `dev|sit|prod/` 多文件，命令行传入根目录即可自动合并，`base/spec.yaml` 提供矩阵/参数模板/工作负载 DSL。
2. 运行命令：
   ```bash
   go run main.go up --env dev --config ./example/config
   ```
3. 如需指定单一云或区域，可附加 `--cloud aws --region ap-northeast-1`，CLI 会读取矩阵展开的目标。
4. 运行完成后 CLI 会输出 JSON，包含 `targets[*].planPreview`、`targets[*].costEstimate`、`targets[*].outputs`、`targets[*].artifacts` 等字段，便于在上层编排中解析。

> 提示：如需自定义配置，可在 `example/config/<env>/` 下新增 YAML 文件，字段会与 `base/` 中的默认值自动合并。

## 2. 触发 GitHub Actions

### 2.1 使用 `gh` CLI
```bash
gh workflow run pulumigo.yaml \
  --repo your-org/your-repo \
  --ref main \
  -f stack_env=prod \
  -f config_path=example/config \
  # 可选: -f stack_cloud=aws -f stack_region=ap-northeast-1
```

### 2.2 使用 REST API
```bash
curl -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <GITHUB_TOKEN>" \
  https://api.github.com/repos/your-org/your-repo/actions/workflows/pulumigo.yaml/dispatches \
  -d '{
    "ref": "main",
    "inputs": {
      "stack_env": "prod",
      "config_path": "example/config"
    }
  }'
```
在工作流中通过 `STACK_ENV`、`CONFIG_PATH`（以及可选的 `STACK_CLOUD`、`STACK_REGION`）环境变量传递到 Runner，解析 CLI 输出即可获得资源规划、成本预估和交付物信息。
如需在 GitHub Workflow inputs 中指定某一云/区域，可额外提供 `stack_cloud`、`stack_region` 字段。

## 3. 触发 GitLab CI

### 3.1 Pipeline Trigger
```bash
curl -X POST \
  -F token=<GITLAB_TRIGGER_TOKEN> \
  -F ref=main \
  -F "variables[STACK_ENV]=sit" \
  -F "variables[CONFIG_PATH]=example/config" \
  # 可选: -F "variables[STACK_CLOUD]=aws" -F "variables[STACK_REGION]=ap-northeast-1" \
  https://gitlab.com/api/v4/projects/<project_id>/trigger/pipeline
```

### 3.2 API v4 Manual Job
若需对已存在的流水线手动触发特定 Job，可调用：
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: <GITLAB_TOKEN>" \
  --form "job=deploy" \
  --form "variables[STACK_ENV]=prod" \
  --form "variables[CONFIG_PATH]=example/config" \
  # 可选: --form "variables[STACK_CLOUD]=aws" --form "variables[STACK_REGION]=ap-northeast-1" \
  https://gitlab.com/api/v4/projects/<project_id>/jobs/<job_id>/play
```

## 4. 解析运行结果

`xcloud up` 的输出形如：
```json
{
  "activeEnv": "dev",
  "configSources": [
    "example/config/base",
    "example/config/dev"
  ],
  "targets": [
    {
      "stack": "dev-aws-ap-northeast-1",
      "cloud": "aws",
      "region": "ap-northeast-1",
      "status": "applied",
      "planPreview": {"create": 5},
      "applied": {"create": 5},
      "costEstimate": {
        "hourlyUSD": 0.05,
        "monthlyUSD": 36.5,
        "breakdownUSD": {"t3.micro": 36.5},
        "assumptions": [
          "730 hours per month",
          "Spot instances estimated at 35% of on-demand pricing"
        ]
      },
      "outputs": {"vpcId": "vpc-0123456789"},
      "artifacts": ["Stack state file: /path/to/Pulumi.dev-aws-ap-northeast-1.yaml"],
      "previewUrl": "https://app.pulumi.com/...",
      "updateUrl": "https://app.pulumi.com/..."
    },
    {
      "stack": "dev-gcp-asia-northeast1",
      "cloud": "gcp",
      "region": "asia-northeast1",
      "status": "skipped",
      "message": "cloud provider gcp not yet supported"
    }
  ]
}
```
CI 平台可直接解析 JSON，针对每个矩阵目标生成资源规划预览、成本预估及交付物报告，或将其回传至上层门户系统。
