// File: src/main.rs

mod config;
mod executor;
mod models;
mod result_store;
mod scheduler;

use crate::config::{load_agent_config, AgentConfig};
use crate::executor::run as run_playbook;
use clap::{Parser, Subcommand};
use std::path::PathBuf;
use tokio::fs;

#[derive(Parser, Debug)]
#[command(name = "cw-agent", version)]
#[command(about = "Xconfig Agent - lightweight local playbook runner")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    /// Run once using playbook(s) from Git repo
    Oneshot,

    /// Run as daemon with interval from config file
    Daemon,

    /// Run full playbook from local file (array of plays)
    Playbook {
        #[arg(short, long)]
        file: PathBuf,
    },

    /// Print latest execution result from local store
    Status,

    /// Show version info
    Version,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let args = Cli::parse();

    // 加载配置文件（共享）
    let agent_config: AgentConfig = load_agent_config("/etc/cw-agent.conf")
        .await
        .unwrap_or_else(|e| {
            eprintln!("⚠️ Failed to load config: {}", e);
            std::process::exit(1);
        });

    match args.command {
        Commands::Oneshot => {
            let repo_dir = "/tmp/cw-agent-sync";
            let branch = agent_config.branch.as_deref().unwrap_or("main");

            config::init_or_update_repo(&agent_config.repo, branch, repo_dir)?;

            for path in &agent_config.playbook {
                let full_path = format!("{}/{}", repo_dir, path);
                let content = fs::read_to_string(&full_path).await?;
                let parsed: Vec<models::Play> = serde_yaml::from_str(&content)?;
                let results = run_playbook(parsed).await?;
                result_store::persist(results).await?;
            }
        }

        Commands::Daemon => {
            scheduler::run_schedule(&agent_config).await?;
        }

        Commands::Playbook { file } => {
            let content = fs::read_to_string(file).await?;
            let parsed: Vec<models::Play> = serde_yaml::from_str(&content)?;
            let results = run_playbook(parsed).await?;
            result_store::persist(results).await?;
        }

        Commands::Status => {
            result_store::print_latest().await?;
        }

        Commands::Version => {
            println!("cw-agent version {}", env!("CARGO_PKG_VERSION"));
        }
    }

    Ok(())
}
