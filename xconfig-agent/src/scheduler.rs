// File: src/scheduler.rs

use crate::config::{check_git_updated, init_or_update_repo, pull_latest, AgentConfig};
use crate::models::Play;
use crate::{executor, result_store};
use std::path::Path;
use tokio::time::{sleep, Duration};

pub async fn run_schedule(agent_config: &AgentConfig) -> anyhow::Result<()> {
    let repo_dir = "/tmp/xconfig-agent-sync";
    let branch = agent_config.branch.as_deref().unwrap_or("main");

    // 启动时 clone 一次
    init_or_update_repo(&agent_config.repo, branch, repo_dir)?;

    let repo_path = Path::new(repo_dir);
    let workdir_prefix = agent_config
        .workdir
        .as_deref()
        .map(Path::new)
        .map(|p| {
            if p.is_absolute() {
                p.to_path_buf()
            } else {
                repo_path.join(p)
            }
        })
        .unwrap_or_else(|| repo_path.to_path_buf());

    loop {
        // 检查是否更新
        if check_git_updated(repo_dir, branch)? {
            println!("🔄 Detected changes in Git repo, updating...");
            pull_latest(repo_dir)?;

            let mut all_results = vec![];

            for path in &agent_config.playbook {
                let playbook_path = workdir_prefix.join(path);
                if playbook_path.exists() {
                    match tokio::fs::read_to_string(&playbook_path).await {
                        Ok(content) => match serde_yaml::from_str::<Vec<Play>>(&content) {
                            Ok(parsed) => match executor::run(parsed).await {
                                Ok(results) => all_results.extend(results),
                                Err(e) => eprintln!("❌ Executor error [{}]: {}", path, e),
                            },
                            Err(e) => eprintln!("❌ YAML parse error [{}]: {}", path, e),
                        },
                        Err(e) => eprintln!("❌ Failed to read file [{}]: {}", path, e),
                    }
                } else {
                    eprintln!("⚠️  Playbook not found: {}", playbook_path.display());
                }
            }

            result_store::persist(all_results).await?;
        } else {
            println!("✅ No changes in Git repo.");
        }

        let interval = agent_config.interval.unwrap_or(60);
        println!("🕒 Sleeping {}s before next check...\n", interval);
        sleep(Duration::from_secs(interval)).await;
    }
}
