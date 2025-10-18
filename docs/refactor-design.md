# Xconfig Refactoring Design

This document outlines the architecture introduced in the recent refactor. The goal is to modularise task execution while providing configurable concurrency and logging.

## Executor
- Encapsulates playbook execution logic in a struct.
- Supports flags:
  - `AggregateOutput` – show aggregated results.
  - `CheckMode` – dry-run without executing tasks.
  - `MaxWorkers` – limit parallel goroutines via a semaphore.
- Optional `LogCollector` interface for custom result handling.
- `MemoryCollector` provides an in-memory implementation used by CLI commands.

## Module Registry
- `internal/modules` hosts built-in modules.
- `registry.go` exposes `Register` and `GetHandler` for dynamic lookup.
- Modules implement the `TaskHandler` type and register themselves via `init`.

## Built-in Modules
- `shell` – run remote shell commands with template rendering support.
- `script` – upload and execute a local script.
- `template` – render a template file on the target host.

## Task Dispatch
- `core/executor/task.go` defines `ExecuteTask` which invokes either a registered module handler or a fallback builtin action.
- Task structs gained a `Type()` helper to determine module names.
- Conditional execution supported with `when` expressions via `EvaluateWhen`.

## CLI Changes
- `remote` and `playbook` commands share the same registry and executor logic.
- `--forks` flag controls concurrency for both commands.
- `remote` runs tasks via `ExecuteTask` with pooled workers and optional logging.

