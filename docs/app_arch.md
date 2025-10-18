# Xconfig Application Architecture

The current implementation focuses on a simple control plane written in Go. Execution happens over SSH by directly running commands on target hosts.

```
+-------------+        +-----------------+        +---------------+
|  xconfig | -----> | SSH connection   | -----> | remote hosts  |
|  CLI        |        | (internal/ssh)   |        | (inventory)   |
+-------------+        +-----------------+        +---------------+
```

1. **Command Parsing** – Cobra commands defined under `cmd/` parse CLI arguments and flags.
2. **Inventory Resolution** – `internal/inventory` reads INI style inventory files and returns host information.
3. **Task Execution** – `core/executor` iterates over plays and tasks. For each task it calls helpers in `internal/ssh` to run shell commands, upload scripts or render templates.
4. **Result Aggregation** – results are printed individually or aggregated using `ssh.AggregatedPrint`.

The architecture is intentionally minimal to ease future expansion:
- new modules like `copy` and `command` will live in the executor and SSH helpers.
- the agent (Rust implementation in `XconfigAgent/`) is planned to run playbooks locally on a schedule.
