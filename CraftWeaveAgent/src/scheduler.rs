// File: src/scheduler.rs

use crate::config::{check_git_updated, init_or_update_repo, pull_latest, AgentConfig};
use crate::models::Play;
use crate::{executor, result_store};
use std::path::Path;
use tokio::time::{sleep, Duration};

pub async fn run_schedule(agent_config: &AgentConfig) -> anyhow::Result<()> {
    let repo_dir = "/tmp/cw-agent-sync";
    let branch = agent_config.branch.as_deref().unwrap_or("main");

    // 启动时 clone 一次
    init_or_update_repo(&agent_config.repo, branch, repo_dir)?;

    loop {
        // 检查是否更新
        if check_git_updated(repo_dir, branch)? {
            println!("🔄 Detected changes in Git repo, updating...");
            pull_latest(repo_dir)?;

            let mut all_results = vec![];

            for path in &agent_config.playbook {
                let full_path = format!("{}/{}", repo_dir, path);
                if Path::new(&full_path).exists() {
                    match tokio::fs::read_to_string(&full_path).await {
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
                    eprintln!("⚠️  Playbook not found: {}", full_path);
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
