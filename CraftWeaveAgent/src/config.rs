// File: src/config.rs

use serde::Deserialize;
use std::fs;
use std::path::Path;
use std::process::Command;
use tokio::fs as tokio_fs;

#[derive(Debug, Deserialize)]
pub struct AgentConfig {
    pub repo: String,
    pub interval: Option<u64>,
    pub playbook: Vec<String>,
    pub branch: Option<String>,
}

/// 加载本地 agent 配置文件（例如 /etc/cw-agent.conf）
pub async fn load_agent_config(path: &str) -> anyhow::Result<AgentConfig> {
    let content = tokio_fs::read_to_string(path).await?;
    let config: AgentConfig = serde_yaml::from_str(&content)?;
    Ok(config)
}

/// 初始 clone 仓库（如果不存在 .git）
pub fn init_or_update_repo(repo: &str, branch: &str, dir: &str) -> anyhow::Result<()> {
    if !Path::new(&format!("{}/.git", dir)).exists() {
        let _ = fs::remove_dir_all(dir);
        fs::create_dir_all(dir)?;
        let status = Command::new("git")
            .args(["clone", "--branch", branch, "--depth", "1", repo, dir])
            .status()?;
        if !status.success() {
            anyhow::bail!("git clone failed");
        }
    }
    Ok(())
}

/// 检查远程仓库是否有更新
pub fn check_git_updated(dir: &str, branch: &str) -> anyhow::Result<bool> {
    let fetch = Command::new("git")
        .current_dir(dir)
        .args(["fetch", "origin"])
        .status()?;
    if !fetch.success() {
        anyhow::bail!("git fetch failed");
    }

    let diff = Command::new("git")
        .current_dir(dir)
        .args(["diff", "--quiet", "HEAD", &format!("origin/{}", branch)])
        .status()?;

    Ok(!diff.success()) // true 表示有变更
}

/// 拉取最新代码（用于更新 playbook）
pub fn pull_latest(dir: &str) -> anyhow::Result<()> {
    let pull = Command::new("git")
        .current_dir(dir)
        .args(["pull", "--rebase"])
        .status()?;
    if !pull.success() {
        anyhow::bail!("git pull failed");
    }
    Ok(())
}
