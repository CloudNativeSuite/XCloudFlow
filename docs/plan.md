# Project Plan

## Overview
This project aims to create a Go-based Ansible-like tool. The initial goal is to establish a simple yet extensible structure for managing remote servers, executing commands, and rendering templates. Additionally, we’ll prepare the codebase for seamless integration with Pulumi’s Go native API.

## Goals
1. Create a robust project directory structure.
2. Implement a flexible command-line interface.
3. Provide placeholders for core modules:
   - Command execution
   - Inventory handling
   - Template rendering
4. Ensure that the design is conducive to future Pulumi integration.

## Steps
1. **Project Initialization**
   - Define the directory layout.
   - Set up `go.mod` and `go.sum` for dependency management.

2. **Command-Line Framework**
   - Integrate Cobra as the CLI framework.
   - Create subcommands for ad-hoc commands (`remote`) and playbooks (`playbook`).
   - Add global flags and configuration file support.

3. **Core Modules**
   - Implement a simple template rendering engine using Go’s `text/template`.
   - Provide an SSH runner for command execution and script upload.
   - Parse INI inventories to obtain target host information.

4. **Pulumi Integration**
   - Outline how Pulumi’s Go SDK will fit into the project.
   - Reserve a directory (`pkg/pulumi/`) for Pulumi-related code.
   - Add a sample function demonstrating a Pulumi resource creation (e.g., an S3 bucket or an EC2 instance).

5. **Documentation**
   - Update this plan as the project evolves.
   - Add usage instructions to `README.md`.

## Current Status
- CLI subcommands `remote` and `playbook` are functional.
- Playbooks support the following task types:
  - `shell` – run a remote command via SSH.
  - `script` – upload and execute a local script.
  - `template` – render a Go template and upload it.
- Roles can be referenced from `roles/<name>/tasks/main.yaml`, but tasks inside the role are limited to the same three module types.
- Aggregated output and dry-run mode are available with `--aggregate` and `--check` flags.

## Future Enhancements
- Dynamic inventory sources from cloud APIs.
- Additional modules (`copy`, `command`, `apt`, `yum`, etc.).
- Variable handling (`set_fact`, `when`, `register`).
- Fully integrated Pulumi workflows for provisioning and configuration.
